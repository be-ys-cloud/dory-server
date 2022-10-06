package tests

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

var baseUrl string

// TestMain contains only the minimum required to start test suite.
func TestMain(m *testing.M) {
	// Generate stack
	pool, network, ldapServer, mail, server, err := setupEnv()
	if err != nil {
		log.Fatalln(err)
	}

	// Set variables that will be used by other test files !
	baseUrl = "http://127.0.0.1:" + server.GetPort("8000/tcp") + "/"

	// Run tests
	_ = m.Run()

	// Destroy stack
	err = destroyEnv(pool, network, ldapServer, mail, server)
	if err != nil {
		log.Fatalln("Could not destroy stack")
	}
}

//--------------- All-in-one methods

func setupEnv() (pool *dockertest.Pool, network *dockertest.Network, ldap *dockertest.Resource, mail *dockertest.Resource, server *dockertest.Resource, err error) {

	pool, err = createPool()
	if err != nil {
		return
	}

	network, err = pool.CreateNetwork("dory-tests")
	if err != nil {
		return
	}

	ldap, err = createOpenLDAPContainer(pool, network)
	if err != nil {
		return
	}

	mail, err = createMailContainer(pool, network)
	if err != nil {
		return
	}

	server, err = createServerContainer(pool, network, ldap.GetPort("636/tcp"), mail.GetPort("1025/tcp"))
	if err != nil {
		return
	}

	// Wait for servers to be up...
	time.Sleep(10 * time.Second)

	return
}

func destroyEnv(pool *dockertest.Pool, network *dockertest.Network, ldap *dockertest.Resource, mail *dockertest.Resource, server *dockertest.Resource) (err error) {

	if ldap != nil {
		if err = deleteContainer(pool, ldap); err != nil {
			return
		}
	}

	if server != nil {
		if err = deleteContainer(pool, server); err != nil {
			return
		}
	}

	if mail != nil {
		if err = deleteContainer(pool, mail); err != nil {
			return
		}
	}

	_ = pool.RemoveNetwork(network)

	path, err := os.Getwd()
	if err != nil {
		return
	}

	_ = exec.Command("rm", "-f", path+"/configuration.json").Run()
	_ = exec.Command("docker", "image", "rm", "-f", "dory_base_test").Run()

	return
}

//--------------- Unit methods

func createPool() (*dockertest.Pool, error) {
	return dockertest.NewPool("")
}

func createOpenLDAPContainer(pool *dockertest.Pool, network *dockertest.Network) (ressource *dockertest.Resource, err error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path = strings.TrimSuffix(path, "/tests")

	ressource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "dory_test_openldap",
		Repository: "osixia/openldap",
		Tag:        "1.5.0",
		Networks:   []*dockertest.Network{network},
		Env: []string{
			"LDAP_TLS_VERIFY_CLIENT=never",
			"LDAP_ORGANISATION=localhost",
			"LDAP_DOMAIN=localhost.priv",
			"LDAP_ADMIN_PASSWORD=admin",
		},
		Mounts: []string{
			path + "/ldap_data/config:/etc/ldap/slapd.d",
			path + "/ldap_data/data:/var/lib/ldap",
		},
	})
	return
}

func createMailContainer(pool *dockertest.Pool, network *dockertest.Network) (ressource *dockertest.Resource, err error) {
	ressource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "dory_test_mailserver",
		Repository: "reachfive/fake-smtp-server",
		Tag:        "latest",
		Networks:   []*dockertest.Network{network},
	})
	return
}

func createServerContainer(pool *dockertest.Pool, network *dockertest.Network, ldapPort string, mailPort string) (*dockertest.Resource, error) {
	ldapPortInt, _ := strconv.Atoi(ldapPort)
	mailPortInt, _ := strconv.Atoi(mailPort)

	configuration := configuration{
		LDAPServer: configurationLdap{
			Admin: configurationLdapAdmin{
				Username: "cn=admin,dc=localhost,dc=priv",
				Password: "admin",
			},
			BaseDN:        "dc=localhost,dc=priv",
			FilterOn:      "(&(objectClass=person)(cn=%s))",
			Address:       "172.17.0.1",
			Port:          ldapPortInt,
			Kind:          "openldap",
			SkipTLSVerify: true,
			EmailField:    "email",
		},
		TOTP: configurationTotp{
			Secret: "AZERTYUIOPQSDFGHJKLMWXCVBN0123456789!",
		},
		MailServer: configurationMail{
			Address:       "172.17.0.1",
			Port:          mailPortInt,
			Password:      "",
			SenderAddress: "noreply@dory.localhost",
			SenderName:    "Dory",
			Subject:       "LDAP Account Management",
			SkipTLSVerify: true,
		},
		FrontAddress: "https://localhost:8001/",
	}

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(configuration)
	err = os.WriteFile(path+"/configuration.json", data, 0777)
	if err != nil {
		return nil, err
	}

	return pool.BuildAndRunWithOptions(strings.TrimSuffix(path, "/test/tests")+"/Dockerfile", &dockertest.RunOptions{
		Name:     "dory_base_test",
		Tag:      "latest",
		Networks: []*dockertest.Network{network},
		Mounts: []string{
			path + "/configuration.json:/app/configuration.json",
		},
	})
}

func deleteContainer(pool *dockertest.Pool, container *dockertest.Resource) error {

	if err := pool.Purge(container); err != nil {
		return err
	}

	return nil
}

// ---- LDAP Configuration

// ---- Structures

type configuration struct {
	LDAPServer   configurationLdap `json:"ldap_server"`
	TOTP         configurationTotp `json:"totp"`
	MailServer   configurationMail `json:"mail_server"`
	FrontAddress string            `json:"front_address"`
}

type configurationLdap struct {
	Admin         configurationLdapAdmin `json:"admin"`
	BaseDN        string                 `json:"base_dn"`
	FilterOn      string                 `json:"filter_on"`
	Address       string                 `json:"address"`
	Port          int                    `json:"port"`
	Kind          string                 `json:"kind"`
	SkipTLSVerify bool                   `json:"skip_tls_verify"`
	EmailField    string                 `json:"email_field"`
}

type configurationLdapAdmin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type configurationTotp struct {
	Secret string `json:"secret"`
}

type configurationMail struct {
	Address       string `json:"address"`
	Port          int    `json:"port"`
	Password      string `json:"password"`
	SenderAddress string `json:"sender_address"`
	SenderName    string `json:"sender_name"`
	Subject       string `json:"subject"`
	SkipTLSVerify bool   `json:"skip_tls_verify"`
}

type email struct {
	Text       string    `json:"text"`
	TextAsHtml string    `json:"textAsHtml"`
	Subject    string    `json:"subject"`
	Date       time.Time `json:"date"`
}

type user struct {
	Username       string         `json:"username"`
	Password       string         `json:"password"`
	TOTP           string         `json:"totp"`
	NewPassword    string         `json:"new_password"`
	OldPassword    string         `json:"old_password"`
	Authentication authentication `json:"authentication"`
}

type totpStruct struct {
	TOTP string `json:"totp"`
}

type authentication struct {
	Token string `json:"token"`
	TOTP  string `json:"totp"`
}
