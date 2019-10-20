package inventory

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	cp "github.com/redhat-cop/dash/pkg/copy"
)

// Process is used to process an inventory found at the provided inventory path/dash.yaml
func (i *Inventory) Process(ns *string) error {

	if i.Namespace != "" {
		ns = &i.Namespace
	}

	log.Println("Namespace is " + *ns)

	if i.Version == 0 {
		log.Println("Unable to determine version. Defaulting to v2...")
		i.Version = 2
	}

	switch i.Version {
	case 3:
		if i.ResourceGroups != nil {
			for _, rg := range i.ResourceGroups {
				err := rg.ProcessResourceGroup(i.Prefix, ns)
				if err != nil {
					return err
				}
			}
		}
	case 2:
		if i.ClusterContentList != nil {
			for _, occ := range i.ClusterContentList {
				err := occ.ProcessClusterContentObject(i.Prefix, ns)
				if err != nil {
					return err
				}
			}
		}
	default:
		log.Fatalln("Unable to determine version. Exiting...")
		os.Exit(1)
	}

	return nil
}

// ProcessClusterContentObject is used to process a v2 format inventory
func (cco *ClusterContentObject) ProcessClusterContentObject(prefix string, ns *string) error {

	file, err := ioutil.TempDir("", "dash")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(file)

	if cco.Content != nil {
		for _, c := range cco.Content {
			err := c.ProcessContent(ns, file, prefix)
			if err != nil {
				return err
			}
		}
	}

	err = Reconcile(file, ns)
	if err != nil {
		return err
	}

	return nil
}

// ProcessResourceGroup is used to process a v3 format inventory
func (rg *ResourceGroup) ProcessResourceGroup(prefix string, ns *string) error {

	if rg.Namespace != "" {
		ns = &rg.Namespace
	}

	file, err := ioutil.TempDir("", "dash")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(file)

	if rg.Resources != nil {
		for _, r := range rg.Resources {
			err := r.ProcessResource(ns, file, prefix)
			if err != nil {
				return err
			}
		}
	}

	err = Reconcile(file, ns)
	if err != nil {
		return err
	}

	return nil
}

// ProcessContent is used to process a single item within a v2 format inventory by converting that item into a v3 format
func (c *ClusterContent) ProcessContent(ns *string, f string, p string) error {
	var CTR Resource

	if c.Namespace != "" {
		ns = &c.Namespace
	}

	log.Println("Resource: " + c.Name + ", Namespace: " + *ns)

	// Converts a v2 ClusterContent to a v3 Resource
	// TODO: Convert Template once it has been added to Resource
	CTR.Name = c.Name
	CTR.Namespace = *ns
	CTR.File = c.File

	if CTR.File != "" {
		err := CTR.ProcessFile(ns, f, p)
		if err != nil {
			return err
		}
	}
	return nil
}

// ProcessResource is used to process a single item within a v3 format inventory
func (r *Resource) ProcessResource(ns *string, f string, p string) error {

	if r.Namespace != "" {
		ns = &r.Namespace
	}

	log.Println("Resource: " + r.Name + ", Namespace: " + *ns)

	// TODO: Determine the type of resource
	if r.File != "" {
		err := r.ProcessFile(ns, f, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func Reconcile(path string, ns *string) error {

	p := path
	abs, err := filepath.Abs(p)
	if err != nil {
		return err
	}
	cmdArgs := []string{"apply", "-f", filepath.Clean(abs)}
	if *ns != "" {
		cmdArgs = append(cmdArgs, "-n", *ns)
	}

	cmd := exec.Command("kubectl", cmdArgs...)
	fmt.Printf("Running command: %s\n", cmd.Args)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n", stdoutStderr)
		return err
	}
	fmt.Printf("%s\n", stdoutStderr)

	return nil

}

func copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	log.Printf("Source: %s, Dest: %s\n", src, dst)
	if sourceFileStat.IsDir() {
		err = cp.CopyDir(src, dst)
		if err != nil {
			return err
		}
		return nil
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}
	err = cp.CopyFile(src, dst)
	if err != nil {
		return err
	}
	return nil
}
