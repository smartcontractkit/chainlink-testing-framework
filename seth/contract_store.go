package seth

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	ErrOpenABIFile = "failed to open ABI file"
	ErrParseABI    = "failed to parse ABI file"
	ErrOpenBINFile = "failed to open BIN file"
	ErrNoABIInFile = "no ABI content found in file"
)

// ContractStore contains all ABIs that are used in decoding. It might also contain contract bytecode for deployment
type ContractStore struct {
	ABIs ABIStore
	BINs map[string][]byte
	mu   *sync.RWMutex
}

type ABIStore map[string]abi.ABI

func (c *ContractStore) GetABI(name string) (*abi.ABI, bool) {
	if !strings.HasSuffix(name, ".abi") {
		name = name + ".abi"
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	abi, ok := c.ABIs[name]
	return &abi, ok
}

func (c *ContractStore) GetAllABIs() []*abi.ABI {
	c.mu.Lock()
	defer c.mu.Unlock()

	var allABIs []*abi.ABI
	for _, a := range c.ABIs {
		aCopy := a //nolint
		allABIs = append(allABIs, &aCopy)
	}

	return allABIs
}

func (c *ContractStore) AddABI(name string, abi abi.ABI) {
	if !strings.HasSuffix(name, ".abi") {
		name = name + ".abi"
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.ABIs[name] = abi
}

func (c *ContractStore) GetBIN(name string) ([]byte, bool) {
	if !strings.HasSuffix(name, ".bin") {
		name = name + ".bin"
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	bin, ok := c.BINs[name]
	return bin, ok
}

func (c *ContractStore) AddBIN(name string, bin []byte) {
	if !strings.HasSuffix(name, ".bin") {
		name = name + ".bin"
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.BINs[name] = bin
}

// NewContractStore creates a new Contract store
func NewContractStore(abiPath, binPath string, gethWrappersPaths []string) (*ContractStore, error) {
	cs := &ContractStore{ABIs: make(ABIStore), BINs: make(map[string][]byte), mu: &sync.RWMutex{}}

	if len(gethWrappersPaths) > 0 && abiPath != "" {
		L.Debug().Msg("ABI files are loaded from both ABI path and Geth wrappers path. This might result in ABI duplication. It shouldn't cause any issues, but it's best to chose only one method.")
	}

	err := cs.loadABIs(abiPath)
	if err != nil {
		return nil, err
	}

	err = cs.loadBINs(binPath)
	if err != nil {
		return nil, err
	}

	err = cs.loadGethWrappers(gethWrappersPaths)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load geth wrappers from %v", gethWrappersPaths)
	}

	return cs, nil
}

func (c *ContractStore) loadABIs(abiPath string) error {
	if abiPath != "" {
		files, err := os.ReadDir(abiPath)
		if err != nil {
			return err
		}
		var foundABI bool
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".abi") {
				L.Debug().Str("File", f.Name()).Msg("ABI file loaded")
				ff, err := os.Open(filepath.Join(abiPath, f.Name()))
				if err != nil {
					return errors.Wrap(err, ErrOpenABIFile)
				}
				a, err := abi.JSON(ff)
				if err != nil {
					return errors.Wrap(err, ErrParseABI)
				}
				c.ABIs[f.Name()] = a
				foundABI = true
			}
		}
		if !foundABI {
			return fmt.Errorf("no ABI files found in '%s'. Fix the path or comment out 'abi_dir' setting", abiPath)
		}
	}

	return nil
}

func (c *ContractStore) loadBINs(binPath string) error {
	if binPath != "" {
		files, err := os.ReadDir(binPath)
		if err != nil {
			return err
		}
		var foundBIN bool
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".bin") {
				L.Debug().Str("File", f.Name()).Msg("BIN file loaded")
				bin, err := os.ReadFile(filepath.Join(binPath, f.Name()))
				if err != nil {
					return errors.Wrap(err, ErrOpenBINFile)
				}
				c.BINs[f.Name()] = common.FromHex(string(bin))
				foundBIN = true
			}
		}
		if !foundBIN {
			return fmt.Errorf("no BIN files found in '%s'. Fix the path or comment out 'bin_dir' setting", binPath)
		}
	}

	return nil
}

func (c *ContractStore) loadGethWrappers(gethWrappersPaths []string) error {
	foundWrappers := false
	for _, gethWrappersPath := range gethWrappersPaths {
		err := filepath.Walk(gethWrappersPath, func(path string, _ os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(path) == ".go" {
				contractName, abiContent, err := extractABIFromGethWrapperDir(path)
				if err != nil {
					if !strings.Contains(err.Error(), ErrNoABIInFile) {
						return err
					}
					L.Debug().Msgf("ABI not found in file due to: %s. Skipping", err.Error())

					return nil
				}
				c.AddABI(contractName, *abiContent)

				// we want to know whether we found at least one wrapper
				if !foundWrappers {
					foundWrappers = true
				}
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	if len(gethWrappersPaths) > 0 && !foundWrappers {
		return fmt.Errorf("no geth wrappers found in '%v'. Fix the path or comment out 'geth_wrappers_dirs' setting", gethWrappersPaths)
	}

	return nil
}

// extractABIFromGethWrapperDir extracts ABI from gethwrappers in a given directory
func extractABIFromGethWrapperDir(filePath string) (string, *abi.ABI, error) {
	fileset := token.NewFileSet()
	node, err := parser.ParseFile(fileset, filePath, nil, parser.AllErrors)
	if err != nil {
		return "", nil, err
	}

	var abiContent string
	// use package name as contract name
	contractName := node.Name.Name

TOP_LOOP:
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		// Loop through the specs (each spec represents a variable or constant declaration)
		for _, spec := range genDecl.Specs {
			vspec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			abiContent = extractValueFromCompositeLiteralField(vspec, "bind", "MetaData", "ABI")
			if abiContent != "" {
				break TOP_LOOP
			}
		}
	}

	if abiContent == "" {
		return "", nil, fmt.Errorf("%s: %s", ErrNoABIInFile, filePath)
	}

	// this cleans up all escape and similar characters that might interfere with the JSON unmarshalling
	var rawAbi interface{}
	if err := json.Unmarshal([]byte(abiContent), &rawAbi); err != nil {
		return "", nil, errors.Wrap(err, "failed to unmarshal ABI content")
	}

	parsedAbi, err := abi.JSON(strings.NewReader(fmt.Sprint(rawAbi)))
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to parse ABI content")
	}

	return contractName, &parsedAbi, nil
}

// extractValueFromCompositeLiteralField finds a composite literal in a given ValueSpec with given type (packageName.typeName)
// and extracts value of a field with a given name
func extractValueFromCompositeLiteralField(vspec *ast.ValueSpec, varPackageName, varType, fieldName string) string {
	for i := range vspec.Names {
		// defensive programming - make sure that for given name index there's a value
		if len(vspec.Values)-1 >= i {
			// check for expected types until we find a field with bind.MetaData type
			// this might need to be updated if the structure of the MetaData struct changes
			// or if package name that stores MetaData changes
			if unaryExpr, ok := vspec.Values[i].(*ast.UnaryExpr); ok {
				if compLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
					if expr, ok := compLit.Type.(*ast.SelectorExpr); ok {
						if x, ok := expr.X.(*ast.Ident); ok {
							if x.Name == varPackageName && expr.Sel.Name == varType {
								return extractStringKeyFromCompositeLiteral(compLit, fieldName)
							}
						}
					}
				}
			}
		}
	}

	return ""
}

// extractStringKeyFromCompositeLiteral returns value of a string field with a given name from a composite literal
func extractStringKeyFromCompositeLiteral(compositeLiteral *ast.CompositeLit, keyName string) string {
	var abiContent string
	for _, elt := range compositeLiteral.Elts {
		if kvExpr, ok := elt.(*ast.KeyValueExpr); ok {
			// Look for filed named "ABI"
			// in a similar way we could extract bytecode from "BIN" field
			if key, ok := kvExpr.Key.(*ast.Ident); ok && key.Name == keyName {
				if abiValue, ok := kvExpr.Value.(*ast.BasicLit); ok && abiValue.Kind == token.STRING {
					abiContent = abiValue.Value
					break
				}
			}
		}
	}

	return abiContent
}
