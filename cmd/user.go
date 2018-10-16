package cmd

import "os"

type User struct {
	Name        string
	Secret      string
	Token       string
	ClusterName string
	Server      string
	Port        string
	Cert        string
}

func (u *User) setUserName(username string) {
	u.Name = username
}

func (u *User) setUserSecret(secret string) {
	u.Secret = secret
}

func (u *User) setUserToken(token []byte) {
	u.Token = string(token)
}

func (u *User) getServerHost() string {
	return "https://" + u.Server + ":" + u.Port
}

func (u *User) getContent() string {
	return u.Name + "-context"
}

func (u *User) getCluster() string {
	return u.ClusterName + "-cluster"
}

func (u *User) CreateStruct() kube {
	user := kube{
		APIVersion: "v1",
		Kind:       "Config",
		Users: []users{
			users{
				Name: u.Name,
				User: user{
					Token: u.Token,
				},
			},
		},
		Clusters: []clusters{
			clusters{
				Name: u.ClusterName,
				Cluster: cluster{
					Server: u.getServerHost(),
					CertificateAuthorityData: u.Cert,
				},
			},
		},
		Contexts: []contexts{
			contexts{
				Name: u.getContent(),
				Context: context{
					Cluster: u.getCluster(),
					User:    u.Name,
				},
			},
		},
		CurrentContext: u.getContent(),
	}
	return user
}

func (u *User) saveUserConfig(data []byte) {
	if _, err := os.Stat(u.ClusterName); os.IsNotExist(err) {
		os.Mkdir(u.ClusterName, 0744)
	}
	f, err := os.Create(u.ClusterName + "/" + u.Name + ".yaml")
	if err != nil {
		errorf("Cannot create file: %v", err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		errorf("Cannot write to file: %v", err)
	}
}
