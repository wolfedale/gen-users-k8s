package cmd

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type kube struct {
	APIVersion     string     `yaml:"apiVersion"`
	Kind           string     `yaml:"kind"`
	Users          []users    `yaml:"users"`
	Clusters       []clusters `yaml:"clusters"`
	Contexts       []contexts `yaml:"contexts"`
	CurrentContext string     `yaml:"current-context"`
}

type users struct {
	Name string `yaml:"name"`
	User user   `yaml:"user"`
}

type user struct {
	Token string `yaml:"token"`
}

type clusters struct {
	Name    string  `yaml:"name"`
	Cluster cluster `yaml:"cluster"`
}

type cluster struct {
	Server                   string `yaml:"server"`
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
}

type contexts struct {
	Context context `yaml:"context"`
	Name    string  `yaml:"name"`
}

type context struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen -kubeconfig $KUBECONFIG",
		Short: "Generating user credentials",
		Long: `
		This tool is calling kubernetes API and generating
		user config files according to their secrets.`,
	}
	return cmd
}

func Run() error {
	// Read toml config file
	c := ReadToml()

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String(
			"kubeconfig",
			filepath.Join(home, ".kube", os.Args[2]),
			"(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String(
			"kubeconfig",
			"",
			"absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// empty struct
	uObj := User{}
	// get the name of our environment name [prod, stage, dev]
	clustername := ClusterName(os.Args[2])

	// add data to the struct
	uObj.ClusterName = clustername
	uObj.Server = c.getHost(clustername)
	uObj.Port = c.getPort(clustername)
	uObj.Cert = c.getCert(clustername)
	users, _ := clientset.CoreV1().ServiceAccounts("kube-system").List(metav1.ListOptions{LabelSelector: "type=users"})
	for _, u := range users.Items {
		uObj.setUserName(u.GetName())
		uObj.setUserSecret(u.Secrets[0].Name)
		s, err := clientset.CoreV1().Secrets("kube-system").Get(u.Secrets[0].Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		uObj.setUserToken(s.Data["token"])
		uObj.saveUserConfig(genYamlFile(uObj.CreateStruct()))
	}
	return nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func genYamlFile(k kube) []byte {
	y, err := yaml.Marshal(k)
	if err != nil {
		errorf("Cannot Unmarshal yaml file: %v", err)
	}
	return y
}
