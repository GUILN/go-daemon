package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

type ConfigData struct {
	Label       string
	ProgramPath string
	ConfigPath  string
	StdOutPath  string
	StdErrPath  string
	KeepAlive   bool
	RunAtLoad   bool
}

var configurationData ConfigData = ConfigData{
	Label:       "com.guiln.observer",
	ProgramPath: "/usr/local/bin/observer",
	ConfigPath:  "/etc/observer.conf",
	StdOutPath:  "/tmp",
	StdErrPath:  "/tmp",
	RunAtLoad:   true,
	KeepAlive:   true,
}

func main() {
	fmt.Printf("Do you like to proceed with the installation of the %s?\n(Y/n)\n", configurationData.Label)
	var usersAnswer string
	fmt.Scanln(&usersAnswer)
	if strings.TrimSpace(strings.ToUpper(usersAnswer)) != "Y" {
		return
	}
	fmt.Printf("Generating Plist file...")
	generatePlistFile()
	fmt.Printf("Copying binary files...")
	if err := copyFile("./observer", configurationData.ProgramPath); err != nil {
		log.Fatalf("Error copying binaries: %v", err)
		return
	}
	fmt.Printf("Copying config files...")
	if err := copyFile("./observer.conf", configurationData.ConfigPath); err != nil {
		log.Fatalf("Error copying configs: %v", err)
		return
	}
	fmt.Printf("Installation succeeded!")
}

func copyFile(source, destination string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	const permissionMask = 0644
	err = ioutil.WriteFile(destination, input, permissionMask)
	if err != nil {
		return err
	}
	return nil
}

func generatePlistFile() {
	plistPath := fmt.Sprintf("%s/Library/LaunchAgents/%s.plist", os.Getenv("HOME"), configurationData.Label)
	f, err := os.Create(plistPath)
	defer f.Close()
	if err != nil {
		log.Fatalf("Template file creation failed: %s", err)
		return
	}

	t := template.Must(template.New("launchdConfig").Parse(plist_template()))
	err = t.Execute(f, configurationData)
	if err != nil {
		log.Fatalf("Template generation failed: %s", err)
	}

}

func plist_template() string {
	return `
<?xml version='1.0' encoding='UTF-8'?>
 <!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\" >
 <plist version='1.0'>
   <dict>
     <key>Label</key><string>{{.Label}}</string>
     <key>ProgramArguments</key>
        <array>
          <string>{{.ProgramPath}}</string>
          <string>-config</string>
          <string>{{.ConfigPath}}</string>
        </array>
	 <key>StandardOutPath</key><string>{{.StdOutPath}}/{{.Label}}.out.log</string>
     <key>StandardErrorPath</key><string>{{.StdErrPath}}/{{.Label}}.err.log</string>
     <key>KeepAlive</key><{{.KeepAlive}}/>
     <key>RunAtLoad</key><{{.RunAtLoad}}/>
   </dict>
</plist>
	`
}
