package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	baseURL string
	token   string
	node    string
	vmID    int
	memory  int
	cores   int
	disk    string
)

var vm = &cobra.Command{
	Use:   "vm",
	Short: "everything related to virtual machines",
	Long:  `This function of the homelab-cli helps managing VMs`,
}

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new virtual machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if err := createVM(name); err != nil {
			fmt.Printf("Error: failed to create VM: %v\n", err)
			return
		}

		fmt.Println("VM created successfully")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an existing virtual machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vmName := args[0]
		// Example logic for deleting a VM
		fmt.Printf("Deleting VM: %s\n", vmName)
	},
}

var renameCmd = &cobra.Command{
	Use:   "rename [oldName] [newName]",
	Short: "Rename an existing virtual machine",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldName := args[0]
		newName := args[1]
		// Example logic for renaming a VM
		fmt.Printf("Renaming VM from %s to %s\n", oldName, newName)
	},
}

func init() {
	vm.AddCommand(createCmd)
	vm.AddCommand(deleteCmd)
	vm.AddCommand(renameCmd)
	createCmd.Flags().StringVar(&baseURL, "url", "https://your-proxmox-server:8006", "Proxmox API base URL")
	createCmd.Flags().StringVar(&token, "token", "", "Proxmox API token (format: 'user@realm!token=uuid')")
	createCmd.Flags().StringVar(&node, "node", "pve", "Proxmox node name")
	createCmd.Flags().IntVar(&vmID, "vmid", 100, "VM ID")
	createCmd.Flags().IntVar(&memory, "memory", 2048, "Memory size in MB")
	createCmd.Flags().IntVar(&cores, "cores", 2, "Number of CPU cores")
	createCmd.Flags().StringVar(&disk, "disk", "local-lvm:10", "Disk configuration string")
	rootCmd.AddCommand(vm)
}

// Functions
func createVM(name string) error {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu", baseURL, node)

	data := map[string]string{
		"vmid":    strconv.Itoa(vmID),
		"name":    name,
		"memory":  strconv.Itoa(memory),
		"cores":   strconv.Itoa(cores),
		"sockets": "1",
		"scsihw":  "virtio-scsi-pci",
		"scsi0":   disk,
		"net0":    "virtio,bridge=vmbr0",
	}

	payload := &bytes.Buffer{}
	if err := json.NewEncoder(payload).Encode(data); err != nil {
		return fmt.Errorf("encoding payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "PVEAPIToken="+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := insecureClient.Do(req)

	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("Proxmox API response:", string(body))
	return nil
}

// Place this globally or reuse across requests
var insecureClient = &http.Client{
	Timeout: time.Second * 10,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}
