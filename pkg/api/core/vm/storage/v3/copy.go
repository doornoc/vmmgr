package v3

import (
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

type File struct {
	uuid string
}

type Progress struct {
	total int64
	size  int64
}

func (p *Progress) Write(data []byte) (int, error) {
	n := len(data)
	p.size += int64(n)

	return n, nil
}

type SSHHandler struct {
	Auth *storage.Auth
}

func NewSSHHandler(handler StorageHandler) *StorageHandler {
	return &handler
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(key)
}

// copy
// wget https://releases.ubuntu.com/22.04/ubuntu-22.04.1-desktop-amd64.iso -O/dev/null |& grep -o "[0-9]\+%"

//func (h *StorageHandler) SFTPRemoteToSFTPRemote() error {
//	//config := &ssh.ClientConfig{User: auth.User, HostKeyCallback: nil, Auth: []ssh.AuthMethod{ssh.Password(auth.Pass)}}
//	// src sftp
//	log.Println(h.SrcAuth.User)
//	configSrc := &ssh.ClientConfig{
//		User:            h.SrcAuth.User,
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//		Auth:            []ssh.AuthMethod{PublicKeyFile("/home/yonedayuto/.ssh/id_rsa")},
//	}
//	h.SrcAuth.IP = "localhost"
//	h.DstAuth.IP = "localhost"
//
//	log.Println(h.SrcAuth.IP + ":22")
//
//	configSrc.SetDefaults()
//	log.Println(h.SrcAuth.IP + ":" + strconv.Itoa(int(h.SrcAuth.Port)))
//	sshSrcConn, err := ssh.Dial("tcp", h.SrcAuth.IP+":"+strconv.Itoa(int(h.SrcAuth.Port)), configSrc)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	defer sshSrcConn.Close()
//
//	// dst sftp
//	configDst := &ssh.ClientConfig{
//		User:            h.DstAuth.User,
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//		//Auth:            []ssh.AuthMethod{ssh.Password(h.Auth.Pass)},
//		Auth: []ssh.AuthMethod{PublicKeyFile("/home/yonedayuto/.ssh/id_rsa")},
//	}
//	configDst.SetDefaults()
//	sshDstConn, err := ssh.Dial("tcp", h.DstAuth.IP+":22", configDst)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	defer sshDstConn.Close()
//
//	// SFTP Src Client
//	srcClient, err := sftp.NewClient(sshSrcConn)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	defer srcClient.Close()
//
//	// SFTP Dst Client
//	dstClient, err := sftp.NewClient(sshDstConn)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	defer dstClient.Close()
//
//	// dstFileの作成
//	dstFile, err := dstClient.Create(h.DstPath)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	defer dstFile.Close()
//
//	// srcFileをOpen
//	srcFile, err := srcClient.Open(h.SrcPath)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//
//	file, err := srcFile.Stat()
//	if err != nil {
//		log.Println("Error: file gateway error")
//		return err
//	}
//
//	log.Println(file.Size())
//
//	p := Progress{total: file.Size()}
//
//	//count := 100
//	//count64 := int64(count)
//	//bar := progressbar.Default(count64)
//
//	// Node側の表示
//	go func() {
//		for {
//			if p.size != p.total {
//				<-time.NewTimer(200 * time.Microsecond).C
//				//bar.Set(int(float64(p.size) / float64(p.total) * 100))
//			} else {
//				return
//			}
//		}
//	}()
//
//	// Node側の表示
//	go func() {
//		for {
//			if p.size != p.total {
//				<-time.NewTimer(1 * time.Second).C
//				//node2.SendServer(h.Input.Info, 0, uint(float64(p.size)/float64(p.total)*100), "Progress: Image Copy", nil)
//			} else {
//				return
//			}
//		}
//	}()
//
//	// コピーの処理
//	bytes, err := io.Copy(dstFile, io.TeeReader(srcFile, &p))
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	//bar.Set(100)
//	fmt.Printf("\n%dbytes copied\n", bytes)
//
//	// sync
//	err = dstFile.Sync()
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	_, err = capacityExpansion(h.DstPath, h.Input.Capacity)
//	if err != nil {
//		log.Println("Error: disk capacity expansion")
//		log.Println(err)
//	}
//
//	//node2.SendServer(h.Input.Info, 0, 100, "Success: Image Copy", nil)
//
//	return nil
//}

//func fileCopy(srcFile, dstFile, controller string) error {
//	log.Println("---Copy disk image")
//	log.Println("src: " + srcFile)
//	log.Println("dst: " + dstFile)
//	src, err := os.Open(srcFile)
//	if err != nil {
//		log.Println("Error: open error")
//		return fmt.Errorf("open error")
//	}
//	defer src.Close()
//	file, err := src.Stat()
//	if err != nil {
//		log.Println("Error: file gateway error")
//		return err
//	}
//
//	dst, err := os.Create(dstFile)
//	if err != nil {
//		log.Println("Error: file create")
//		return err
//	}
//	defer dst.Close()
//
//	p := Progress{total: file.Size()}
//
//	count := 100
//	count64 := int64(count)
//	bar := progressbar.Default(count64)
//
//	go func() {
//		for {
//			if p.size != p.total {
//				<-time.NewTimer(200 * time.Microsecond).C
//				//log.Println(tmp.fileSize)
//				bar.Set(int(float64(p.size) / float64(p.total) * 100))
//				//sendServer()
//			} else {
//				log.Println("end")
//				return
//			}
//		}
//	}()
//
//	_, err = io.Copy(dst, io.TeeReader(src, &p))
//	if err != nil {
//		log.Println("Error: file copy error")
//		return err
//	}
//
//	return nil
//}
