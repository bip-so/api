package canvasbranch

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

type BranchIDUUID struct {
	id   uint64
	uuid string
}

func (s canvasBranchService) UtilCreateUUIDArray(m map[string]BranchIDUUID) []string {
	var strArray []string
	for k, _ := range m {
		strArray = append(strArray, k)
	}
	return strArray
}

// Returns a Map of all the UUID : {ID, UUID} for a Branch
func (s canvasBranchService) UtilBranchLookupMapMaker(branchID uint64) map[string]BranchIDUUID {
	// Build Maps of IDs and UUI's to start with
	BlockMapIDUUID := make(map[string]BranchIDUUID)

	// Get List of Blocks for this Branch
	blocks, _ := App.Service.GetAllBlockByBranchID(branchID)
	// Loop through the List
	for _, block := range *blocks {
		BlockMapIDUUID[block.UUID.String()] = BranchIDUUID{
			id:   block.ID,
			uuid: block.UUID.String(),
		}
	}
	return BlockMapIDUUID
}

// deleted := ArrayDiff(parentBranch, roughBranch)
func (s canvasBranchService) UtilBlocksToDelete(roughB map[string]BranchIDUUID, parentB map[string]BranchIDUUID, finalBlocksMap map[string]*models.Block) []string {
	roughBUUIDArray := s.UtilCreateUUIDArray(roughB)
	parentBUUIDArray := s.UtilCreateUUIDArray(parentB)
	diff := s.ArrayDiff(parentBUUIDArray, roughBUUIDArray)
	for _, uuid := range diff {
		if _, exists := finalBlocksMap[uuid]; exists {
			diff = utils.RemoveItemFromSlice(diff, uuid)
		}
	}
	return diff
}

// created := ArrayDiff(roughBranch, parentBranch)
func (s canvasBranchService) UtilBlocksToCreate(roughB map[string]BranchIDUUID, parentB map[string]BranchIDUUID, finalBlocksMap map[string]*models.Block) []string {
	roughBUUIDArray := s.UtilCreateUUIDArray(roughB)
	parentBUUIDArray := s.UtilCreateUUIDArray(parentB)
	diff := s.ArrayDiff(roughBUUIDArray, parentBUUIDArray)
	for _, uuid := range diff {
		if _, exists := finalBlocksMap[uuid]; !exists {
			diff = utils.RemoveItemFromSlice(diff, uuid)
		}
	}
	return diff
}

func (s canvasBranchService) UtilBlocksToUpdate(roughB map[string]BranchIDUUID, parentB map[string]BranchIDUUID) []string {
	roughBUUIDArray := s.UtilCreateUUIDArray(roughB)
	parentBUUIDArray := s.UtilCreateUUIDArray(parentB)
	return s.UtilIntersection(roughBUUIDArray, parentBUUIDArray)
}

// Array Nonsense

func (s canvasBranchService) UtilIntersection(s1, s2 []string) (inter []string) {
	hash := make(map[string]bool)
	for _, e := range s1 {
		hash[e] = true
	}
	for _, e := range s2 {
		// If elements present in the hashmap then append intersection list.
		if hash[e] {
			inter = append(inter, e)
		}
	}
	//Remove dups from slice.
	inter = s.removeDups(inter)
	return
}

//Remove dups from slice.
func (s canvasBranchService) removeDups(elements []string) (nodups []string) {
	encountered := make(map[string]bool)
	for _, element := range elements {
		if !encountered[element] {
			nodups = append(nodups, element)
			encountered[element] = true
		}
	}
	return
}

func (s canvasBranchService) InSlice(n string, h []string) bool {
	for _, v := range h {
		if v == n {
			return true
		}
	}
	return false
}

// ArrayDiff Get the elements of an List1 which are not in List2
// @todo: https://github.com/juliangruber/go-intersect/blob/master/intersect.go
// Tested: https://go.dev/play/p/HXbsa2IUTEc
func (s canvasBranchService) ArrayDiff(list1 []string, list2 []string) []string {
	res := make([]string, 0)
	for _, value := range list1 {
		if !s.InSlice(value, list2) {
			res = append(res, value)
		}
	}
	return res
}

// func (s canvasBranchService) Merg {}
