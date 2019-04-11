package cloudhealth

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nextgenhealthcare/cloudhealth-sdk-go"
)

func resourceCloudHealthPerspective() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudHealthPerspectiveCreate,
		Read:   resourceCloudHealthPerspectiveRead,
		Update: resourceCloudHealthPerspectiveUpdate,
		Delete: resourceCloudHealthPerspectiveDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"include_in_reports": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"group": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"ref_id": {
							Type:     schema.TypeString,
							ForceNew: false,
							Computed: true,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
							Default:  "filter",
						},
						"rule": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"asset": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: false,
									},
									// for type="categorize"
									"tag_field": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: false,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									// for type="categorize"
									"field": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: false,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"combine_with": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: false,
									},
									"condition": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tag_field": {
													Type:     schema.TypeList,
													Optional: true,
													ForceNew: false,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"field": {
													Type:     schema.TypeList,
													Optional: true,
													ForceNew: false,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"op": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: false,
													Default:  "=",
												},
												"val": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: false,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"constant": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"constant_type": {
							Type:     schema.TypeString,
							ForceNew: false,
							Computed: true,
						},
						"ref_id": {
							Type:     schema.TypeString,
							ForceNew: false,
							Computed: true,
						},
						"blk_id": {
							Type:     schema.TypeString,
							ForceNew: false,
							Computed: true,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							ForceNew: false,
							Computed: true,
							Optional: true,
						},
						"val": {
							Type:     schema.TypeString,
							ForceNew: false,
							Computed: true,
							Optional: true,
						},
						"is_other": {
							Type:     schema.TypeString,
							ForceNew: false,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceCloudHealthPerspectiveCreate(d *schema.ResourceData, m interface{}) error {
	var createdId string
	client := m.(*cloudhealth.Client)
	perspective, err := convertPerspective(d)
	if err != nil {
		return fmt.Errorf("Could not convert perspective: %v", err)
	}

	createdId, err = client.CreatePerspective(perspective)
	if err != nil {
		return fmt.Errorf("Could not create perspective: %v", err)
	}

	d.SetId(createdId)
	// We need to set the constants field to what cloudhealth thinks it is, as
	// its computed we need to read it back from cloudhealth - easiest to do that
	// by using the read method
	return resourceCloudHealthPerspectiveRead(d, m)
}

func resourceCloudHealthPerspectiveRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudhealth.Client)

	id := d.Id()
	perspective, err := client.GetPerspective(id)

	switch err {
	case cloudhealth.ErrPerspectiveNotFound:
		d.SetId("")
		return nil
	default:
		return fmt.Errorf("Error when reading perspective %s: %v", id, err)
	}

	return buildPerspective(perspective, d)
}

func resourceCloudHealthPerspectiveUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudhealth.Client)

	id := d.Id()
	perspective, err := convertPerspective(d)
	if err != nil {
		return fmt.Errorf("Could not convert perspective: %v", err)
	}

	_, err = client.UpdatePerspective(id, perspective)
	if err != nil {
		return fmt.Errorf("Could not create perspective: %v", err)
	}

	// We need to set the constants field to what cloudhealth thinks it is, as
	// its computed we need to read it back from cloudhealth - easiest to do that
	// by using the read method
	return resourceCloudHealthPerspectiveRead(d, m)
}

func resourceCloudHealthPerspectiveDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudhealth.Client)
	err := client.DeletePerspective(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func convertPerspective(d *schema.ResourceData) (perspective *cloudhealth.Perspective, err error) {
	constants := []*cloudhealth.Constant{
		cloudhealth.NewConstant(cloudhealth.StaticGroupType),
		cloudhealth.NewConstant(cloudhealth.DynamicGroupType),
		cloudhealth.NewConstant(cloudhealth.DynamicGroupBlockType),
	}

	constantsByType := make(map[string]*cloudhealth.Constant)
	for _, constant := range constants {
		constantsByType[constant.Type] = constant
	}

	name, ok := d.GetOk("name")
	if !ok {
		return nil, fmt.Errorf("Required name field")
	}
	includeInReports := d.Get("include_in_reports")

	perspective = new(cloudhealth.Perspective)
	perspective.Schema.Name = name.(string)
	perspective.Schema.IncludeInReports = strconv.FormatBool(includeInReports.(bool))

	tfGroups := getArray(d, "group")
	tfConstants := getArray(d, "constant")

	if len(tfGroups) > 0 {
		err = fixRefIDs(tfGroups, tfConstants)
		if err != nil {
			return nil, err
		}

		err = d.Set("group", tfGroups)
		if err != nil {
			return nil, err
		}
	}

	for _, tfGroup := range tfGroups {
		tfGroup := tfGroup.(map[string]interface{})
		refId := tfGroup["ref_id"].(string)
		name := tfGroup["name"].(string)
		groupType := tfGroup["type"].(string)

		var constantType string

		if tfGroup["type"].(string) == "categorize" {
			// Convert any dynamic groups for this group (if it's a Dynamic Group Block)
			dynamicGroupConstantItems := convertDynamicGroupConstantItems(refId, tfConstants)
			constantsByType[cloudhealth.DynamicGroupType].List = append(constantsByType[cloudhealth.DynamicGroupType].List, dynamicGroupConstantItems...)
			constantType = cloudhealth.DynamicGroupBlockType
		} else if tfGroup["type"].(string) == "filter" {
			constantType = cloudhealth.StaticGroupType
		} else {
			return nil, fmt.Errorf("Unknown group type: %s. Expected filter or categorize", tfGroup["type"])
		}

		// Convert any rules
		rules, err := convertRules(refId, name, groupType, tfGroup["rule"].([]interface{}))
		if err != nil {
			return nil, err
		}
		perspective.Schema.Rules = append(perspective.Schema.Rules, rules...)

		// Add a constant for this group
		constantItem := cloudhealth.ConstantItem{
			Name:  name,
			RefID: refId,
		}
		constant := constantsByType[constantType]
		constant.List = append(constant.List, constantItem)
	}

	err = addOtherConstants(tfConstants, constantsByType)
	if err != nil {
		return nil, err
	}

	// Only add constants that have something in them
	for _, constantGroup := range constants {
		if len(constantGroup.List) > 0 {
			perspective.Schema.Constants = append(perspective.Schema.Constants, *constantGroup)
		}
	}
	perspective.Schema.Merges = make([]interface{}, 0)
	return perspective, nil
}

func fixRefIDs(groups []interface{}, constants []interface{}) error {
	/* This is to reconcile the ref_id on groups with the ones in constants.

	   Groups are an ordered list and yet also identified by their ref_id.
	   There's no direct way to express an ordered map in terraform schema, so we
	   use a list. When groups are reordered, the computed ref_id fields stay put;
	   they do not follow the rest of the groups contents.

	   So we use the "constants" structure to reconcile these situations.

	   If the group is renamed in-place, the new name won't have an entry in
	   constants, so it's presumed to keep its ref_id.

	   If the group is re-ordered, we look up the ref_ids by the name in the
	   constants structure.

	   If you both re-order and re-name, it will correct the ref_ids of all the
	   other groups, but the reordered group with the new name will be given a
	   new ref_id
	*/

	refIdByNameFromConstants := make(map[string]string)
	maxRefId := 0
	for _, c := range constants {
		c := c.(map[string]interface{})
		refIdByNameFromConstants[c["name"].(string)] = c["ref_id"].(string)
		constantRefIdInt, err := strconv.Atoi(c["ref_id"].(string))
		if err != nil {
			return fmt.Errorf("Group with non integer ref_id: %s", c["ref_id"])
		}
		if constantRefIdInt >= maxRefId {
			maxRefId = constantRefIdInt + 1
		}
	}
	usedRefIds := make(map[string]bool)

	// Go through and apply the ref_id from the constant to anything that matches the same name in the group
	for _, g := range groups {
		g := g.(map[string]interface{})
		groupName := g["name"].(string)
		if constantRefId, ok := refIdByNameFromConstants[groupName]; ok {
			g["ref_id"] = constantRefId
			if usedRefIds[constantRefId] == true {
				return fmt.Errorf("Two groups with the same name: %s", groupName)
			}
			usedRefIds[constantRefId] = true
		}
	}

	// Now for any group who name is not in constants, either use its exising
	// ref_id (we assume this meant a rename) or, if it doesn't have one,
	// generate a unique one
	for _, g := range groups {
		g := g.(map[string]interface{})
		groupName := g["name"].(string)
		if _, inConstants := refIdByNameFromConstants[groupName]; inConstants {
			// Already fixed ref_id above
			continue
		}

		groupRefId := g["ref_id"].(string)
		if groupRefId != "" && usedRefIds[groupRefId] == false {
			// Group was renamed; stick with the existing groupRefId
			continue
		}

		// Group is new - assign a new ref id
		// Must be an integer that is not already in use
		g["ref_id"] = strconv.Itoa(maxRefId)
		maxRefId++
	}

	return nil
}

func convertDynamicGroupConstantItems(groupRefId string, constants []interface{}) []cloudhealth.ConstantItem {
	result := make([]cloudhealth.ConstantItem, 0)

	for _, c := range constants {
		c := c.(map[string]interface{})
		if c["blk_id"] != groupRefId {
			continue
		}
		blk_id := groupRefId
		result = append(result, cloudhealth.ConstantItem{
			Name:  c["name"].(string),
			RefID: c["ref_id"].(string),
			BlkID: &blk_id,
			Val:   c["val"].(string),
		})
	}
	return result
}

func convertRules(groupRefId string, groupName string, groupType string, rules []interface{}) (result []cloudhealth.Rule, err error) {
	result = make([]cloudhealth.Rule, len(rules))

	for ruleIdx, r := range rules {
		r := r.(map[string]interface{})

		rj := &result[ruleIdx]

		rj.Type = groupType
		if groupType == "categorize" {
			rj.RefID = groupRefId
			rj.Name = groupName
		} else if groupType == "filter" {
			rj.To = groupRefId
		} else {
			return nil, fmt.Errorf("Unrecognized group type %s", groupType)
		}

		rj.Asset = stringOrNil(r["asset"])
		rj.Field = convertStringArray(r["field"])
		rj.TagField = convertStringArray(r["tag_field"])

		if r["condition"] != nil {
			rj.Condition = convertConditions(r["condition"].([]interface{}), stringOrNil(r["combine_with"]))
		} else {
			rj.Condition = nil
		}
	}
	return result, nil
}

func convertConditions(conditions []interface{}, combineWith string) (result *cloudhealth.Condition) {
	if len(conditions) == 0 {
		return nil
	}
	result = new(cloudhealth.Condition)
	result.Clauses = make([]cloudhealth.Clause, len(conditions))
	result.CombineWith = combineWith
	for idx, condition := range conditions {
		condition := condition.(map[string]interface{})
		result.Clauses[idx] = cloudhealth.Clause{
			Field:    convertStringArray(condition["field"]),
			TagField: convertStringArray(condition["tag_field"]),
			Op:       stringOrNil(condition["op"]),
			Val:      stringOrNil(condition["val"]),
		}
	}
	return result
}

func convertConstant(tfConstant map[string]interface{}) (constantType string, constantItem cloudhealth.ConstantItem) {
	constantType = tfConstant["constant_type"].(string)
	constantItem = cloudhealth.ConstantItem{
		RefID: stringOrNil(tfConstant["ref_id"]),
		Name:  stringOrNil(tfConstant["name"]),
		Val:   stringOrNil(tfConstant["val"]),
	}
	if constantType == cloudhealth.DynamicGroupType {
		blk_id := stringOrNil(tfConstant["blk_id"])
		constantItem.BlkID = &blk_id
	}
	if stringOrNil(tfConstant["is_other"]) == "true" {
		constantItem.IsOther = "true"
	}
	return constantType, constantItem
}

func addOtherConstants(tfConstants []interface{}, constantsByType map[string]*cloudhealth.Constant) error {
	// Add "other" constants
	// These are constants that have literally is_other == "true" or dynamic
	// groups with empty blk_ids
	for _, tfConstant := range tfConstants {
		tfConstant := tfConstant.(map[string]interface{})

		if tfConstant["is_other"].(string) == "true" ||
			(tfConstant["constant_type"].(string) == cloudhealth.DynamicGroupType && tfConstant["blk_id"] == "") {

			constantType, constantItem := convertConstant(tfConstant)

			constant := constantsByType[constantType]
			if constant == nil {
				return fmt.Errorf("Unknown constant type %s", constantType)
			}
			constant.List = append(constant.List, constantItem)
		}
	}
	return nil
}

func convertStringArray(maybeStringArray interface{}) []string {
	if maybeStringArray == nil {
		return nil
	}
	ss := maybeStringArray.([]interface{})
	result := make([]string, len(ss))
	for idx, s := range ss {
		result[idx] = s.(string)
	}
	return result
}

func stringOrNil(s interface{}) string {
	if s == nil {
		return ""
	}
	return s.(string)
}

func getArray(d *schema.ResourceData, field string) []interface{} {
	if v, ok := d.GetOk(field); ok {
		return v.([]interface{})
	} else {
		return make([]interface{}, 0)
	}
}

func buildPerspective(p *cloudhealth.Perspective, d *schema.ResourceData) error {
	d.Set("name", p.Schema.Name)

	if v, err := strconv.ParseBool(p.Schema.IncludeInReports); err == nil {
		err = d.Set("include_in_reports", v)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	groupByRef := buildGroups(p)
	groups, err := populateRules(p, groupByRef)
	if err != nil {
		return err
	}

	constants := buildConstants(p)

	d.Set("group", groups)

	err = d.Set("constant", constants)
	if err != nil {
		return err
	}
	return nil
}

func buildGroups(p *cloudhealth.Perspective) (groupByRef map[string]cloudhealth.Group) {
	groupByRef = make(map[string]cloudhealth.Group)

	for _, constant := range p.Schema.Constants {
		if constant.Type != cloudhealth.StaticGroupType && constant.Type != cloudhealth.DynamicGroupBlockType {
			continue
		}
		for _, constantGroup := range constant.List {
			if constantGroup.IsOther == "true" {
				// An "other" group, solely handled by buildConstants()
				continue
			}
			group := make(cloudhealth.Group)
			group["name"] = constantGroup.Name
			group["ref_id"] = constantGroup.RefID
			group["rule"] = make([]map[string]interface{}, 0)
			if constant.Type == cloudhealth.DynamicGroupBlockType {
				group["type"] = "categorize"
			} else {
				group["type"] = "filter"
			}
			groupByRef[constantGroup.RefID] = group
		}
	}
	return groupByRef
}

func populateRules(p *cloudhealth.Perspective, groupByRef map[string]cloudhealth.Group) (groups []cloudhealth.Group, err error) {
	groupByRefSeen := make(map[string]bool)
	groups = make([]cloudhealth.Group, 0)
	for _, jsonRule := range p.Schema.Rules {
		groupRef := jsonRule.To
		if groupRef == "" {
			groupRef = jsonRule.RefID
			if groupRef == "" {
				return nil, fmt.Errorf("Unable to find 'to' for rule for asset %s", jsonRule.Type)
			}
		}

		rule := make(map[string]interface{})
		group := groupByRef[groupRef]
		if group == nil {
			return nil, fmt.Errorf("Group reference %s not found", groupRef)
		}

		// Order the groups by order that the rules are seen.  CHT technically
		// allows the groups for rules to be interleaved, but this is horribly
		// confusing and not in the UI
		if groupByRefSeen[groupRef] == false {
			groups = append(groups, group)
			groupByRefSeen[groupRef] = true
		}

		group["rule"] = append(group["rule"].([]map[string]interface{}), rule)

		if jsonRule.Type != group["type"] {
			return nil, fmt.Errorf("Unknown rule type %s; expected %s", jsonRule.Type, group["type"])
		}
		rule["asset"] = jsonRule.Asset
		if jsonRule.TagField != nil {
			rule["tag_field"] = jsonRule.TagField
		}
		if jsonRule.Field != nil {
			rule["field"] = jsonRule.Field
		}

		if jsonRule.Condition != nil {
			rule["combine_with"] = jsonRule.Condition.CombineWith
			jsonClauses := jsonRule.Condition.Clauses
			if jsonClauses != nil {
				rule["condition"] = buildCondition(jsonClauses)
			}
		}
	}

	return groups, nil
}

func buildCondition(srcClauses []cloudhealth.Clause) (clauses []map[string]interface{}) {
	clauses = make([]map[string]interface{}, len(srcClauses))

	for idx, srcClause := range srcClauses {
		clause := make(map[string]interface{})
		clauses[idx] = clause

		if srcClause.TagField != nil {
			clause["tag_field"] = srcClause.TagField
		}
		if srcClause.Field != nil {
			clause["field"] = srcClause.Field
		}
		clause["op"] = srcClause.Op
		clause["val"] = srcClause.Val
	}
	return clauses
}

func buildConstants(p *cloudhealth.Perspective) []cloudhealth.Group {
	result := make([]cloudhealth.Group, 0)
	for _, srcConstant := range p.Schema.Constants {
		for _, srcConstantGroup := range srcConstant.List {
			constant := cloudhealth.Group{
				"constant_type": srcConstant.Type,
				"ref_id":        srcConstantGroup.RefID,
				"name":          srcConstantGroup.Name,
				"val":           srcConstantGroup.Val,
				"is_other":      srcConstantGroup.IsOther,
			}
			if srcConstantGroup.BlkID != nil {
				constant["blk_id"] = *srcConstantGroup.BlkID
			}

			result = append(result, constant)
		}
	}

	return result
}
