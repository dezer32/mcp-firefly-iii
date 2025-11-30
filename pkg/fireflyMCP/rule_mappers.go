package fireflyMCP

import (
	"github.com/dezer32/mcp-firefly-iii/pkg/client"
)

// mapRuleTriggerToDTO converts client.RuleTrigger to RuleTrigger DTO
func mapRuleTriggerToDTO(trigger *client.RuleTrigger) RuleTrigger {
	if trigger == nil {
		return RuleTrigger{}
	}

	result := RuleTrigger{
		Id:             getStringValue(trigger.Id),
		Type:           string(trigger.Type),
		Value:          trigger.Value,
		Prohibited:     false,
		Active:         true,
		StopProcessing: false,
		Order:          0,
	}

	if trigger.Prohibited != nil {
		result.Prohibited = *trigger.Prohibited
	}
	if trigger.Active != nil {
		result.Active = *trigger.Active
	}
	if trigger.StopProcessing != nil {
		result.StopProcessing = *trigger.StopProcessing
	}
	if trigger.Order != nil {
		result.Order = int(*trigger.Order)
	}

	return result
}

// mapRuleActionToDTO converts client.RuleAction to RuleAction DTO
func mapRuleActionToDTO(action *client.RuleAction) RuleAction {
	if action == nil {
		return RuleAction{}
	}

	result := RuleAction{
		Id:             getStringValue(action.Id),
		Type:           string(action.Type),
		Value:          action.Value,
		Active:         true,
		StopProcessing: false,
		Order:          0,
	}

	if action.Active != nil {
		result.Active = *action.Active
	}
	if action.StopProcessing != nil {
		result.StopProcessing = *action.StopProcessing
	}
	if action.Order != nil {
		result.Order = int(*action.Order)
	}

	return result
}

// mapRuleReadToRule converts client.RuleRead to Rule DTO
func mapRuleReadToRule(ruleRead *client.RuleRead) *Rule {
	if ruleRead == nil {
		return nil
	}

	rule := &Rule{
		Id:             ruleRead.Id,
		Title:          ruleRead.Attributes.Title,
		Description:    ruleRead.Attributes.Description,
		RuleGroupId:    ruleRead.Attributes.RuleGroupId,
		RuleGroupTitle: ruleRead.Attributes.RuleGroupTitle,
		Order:          0,
		Trigger:        string(ruleRead.Attributes.Trigger),
		Active:         true,
		Strict:         true,
		StopProcessing: false,
		Triggers:       []RuleTrigger{},
		Actions:        []RuleAction{},
	}

	if ruleRead.Attributes.Order != nil {
		rule.Order = int(*ruleRead.Attributes.Order)
	}
	if ruleRead.Attributes.Active != nil {
		rule.Active = *ruleRead.Attributes.Active
	}
	if ruleRead.Attributes.Strict != nil {
		rule.Strict = *ruleRead.Attributes.Strict
	}
	if ruleRead.Attributes.StopProcessing != nil {
		rule.StopProcessing = *ruleRead.Attributes.StopProcessing
	}

	// Map triggers
	for _, trigger := range ruleRead.Attributes.Triggers {
		rule.Triggers = append(rule.Triggers, mapRuleTriggerToDTO(&trigger))
	}

	// Map actions
	for _, action := range ruleRead.Attributes.Actions {
		rule.Actions = append(rule.Actions, mapRuleActionToDTO(&action))
	}

	return rule
}

// mapRuleArrayToRuleList converts client.RuleArray to RuleList DTO
func mapRuleArrayToRuleList(ruleArray *client.RuleArray) *RuleList {
	if ruleArray == nil {
		return &RuleList{
			Data: []Rule{},
		}
	}

	ruleList := &RuleList{
		Data: make([]Rule, 0, len(ruleArray.Data)),
	}

	for _, ruleRead := range ruleArray.Data {
		if mappedRule := mapRuleReadToRule(&ruleRead); mappedRule != nil {
			ruleList.Data = append(ruleList.Data, *mappedRule)
		}
	}

	// Map pagination
	if ruleArray.Meta.Pagination != nil {
		pagination := ruleArray.Meta.Pagination
		ruleList.Pagination = Pagination{
			Count:       getIntValue(pagination.Count),
			Total:       getIntValue(pagination.Total),
			CurrentPage: getIntValue(pagination.CurrentPage),
			PerPage:     getIntValue(pagination.PerPage),
			TotalPages:  getIntValue(pagination.TotalPages),
		}
	}

	return ruleList
}

// mapRuleGroupReadToRuleGroup converts client.RuleGroupRead to RuleGroup DTO
func mapRuleGroupReadToRuleGroup(ruleGroupRead *client.RuleGroupRead) *RuleGroup {
	if ruleGroupRead == nil {
		return nil
	}

	ruleGroup := &RuleGroup{
		Id:          ruleGroupRead.Id,
		Title:       ruleGroupRead.Attributes.Title,
		Description: ruleGroupRead.Attributes.Description,
		Order:       0,
		Active:      true,
	}

	if ruleGroupRead.Attributes.Order != nil {
		ruleGroup.Order = int(*ruleGroupRead.Attributes.Order)
	}
	if ruleGroupRead.Attributes.Active != nil {
		ruleGroup.Active = *ruleGroupRead.Attributes.Active
	}

	return ruleGroup
}

// mapRuleGroupArrayToRuleGroupList converts client.RuleGroupArray to RuleGroupList DTO
func mapRuleGroupArrayToRuleGroupList(ruleGroupArray *client.RuleGroupArray) *RuleGroupList {
	if ruleGroupArray == nil {
		return &RuleGroupList{
			Data: []RuleGroup{},
		}
	}

	ruleGroupList := &RuleGroupList{
		Data: make([]RuleGroup, 0, len(ruleGroupArray.Data)),
	}

	for _, ruleGroupRead := range ruleGroupArray.Data {
		if mappedRuleGroup := mapRuleGroupReadToRuleGroup(&ruleGroupRead); mappedRuleGroup != nil {
			ruleGroupList.Data = append(ruleGroupList.Data, *mappedRuleGroup)
		}
	}

	// Map pagination
	if ruleGroupArray.Meta.Pagination != nil {
		pagination := ruleGroupArray.Meta.Pagination
		ruleGroupList.Pagination = Pagination{
			Count:       getIntValue(pagination.Count),
			Total:       getIntValue(pagination.Total),
			CurrentPage: getIntValue(pagination.CurrentPage),
			PerPage:     getIntValue(pagination.PerPage),
			TotalPages:  getIntValue(pagination.TotalPages),
		}
	}

	return ruleGroupList
}

// mapRuleTriggerRequestToStore converts RuleTriggerRequest to client.RuleTriggerStore
func mapRuleTriggerRequestToStore(req RuleTriggerRequest) client.RuleTriggerStore {
	store := client.RuleTriggerStore{
		Type:  client.RuleTriggerKeyword(req.Type),
		Value: req.Value,
	}

	if req.Prohibited != nil {
		store.Prohibited = req.Prohibited
	}
	if req.Active != nil {
		store.Active = req.Active
	}
	if req.StopProcessing != nil {
		store.StopProcessing = req.StopProcessing
	}

	return store
}

// mapRuleActionRequestToStore converts RuleActionRequest to client.RuleActionStore
func mapRuleActionRequestToStore(req RuleActionRequest) client.RuleActionStore {
	store := client.RuleActionStore{
		Type:  client.RuleActionKeyword(req.Type),
		Value: req.Value,
	}

	if req.Active != nil {
		store.Active = req.Active
	}
	if req.StopProcessing != nil {
		store.StopProcessing = req.StopProcessing
	}

	return store
}

// mapRuleStoreRequestToAPI converts RuleStoreRequest to client.RuleStore
func mapRuleStoreRequestToAPI(req *RuleStoreRequest) client.RuleStore {
	store := client.RuleStore{
		Title:          req.Title,
		Description:    req.Description,
		RuleGroupId:    req.RuleGroupId,
		RuleGroupTitle: req.RuleGroupTitle,
		Trigger:        client.RuleTriggerType(req.Trigger),
		Active:         req.Active,
		Strict:         req.Strict,
		StopProcessing: req.StopProcessing,
	}

	// Map triggers
	triggers := make([]client.RuleTriggerStore, len(req.Triggers))
	for i, trigger := range req.Triggers {
		triggers[i] = mapRuleTriggerRequestToStore(trigger)
	}
	store.Triggers = triggers

	// Map actions
	actions := make([]client.RuleActionStore, len(req.Actions))
	for i, action := range req.Actions {
		actions[i] = mapRuleActionRequestToStore(action)
	}
	store.Actions = actions

	return store
}

// mapRuleTriggerRequestToUpdate converts RuleTriggerRequest to client.RuleTriggerUpdate
func mapRuleTriggerRequestToUpdate(req RuleTriggerRequest) client.RuleTriggerUpdate {
	update := client.RuleTriggerUpdate{
		Type:  (*client.RuleTriggerKeyword)(&req.Type),
		Value: &req.Value,
	}

	if req.Active != nil {
		update.Active = req.Active
	}
	if req.StopProcessing != nil {
		update.StopProcessing = req.StopProcessing
	}

	return update
}

// mapRuleActionRequestToUpdate converts RuleActionRequest to client.RuleActionUpdate
func mapRuleActionRequestToUpdate(req RuleActionRequest) client.RuleActionUpdate {
	update := client.RuleActionUpdate{
		Type:  (*client.RuleActionKeyword)(&req.Type),
		Value: req.Value,
	}

	if req.Active != nil {
		update.Active = req.Active
	}
	if req.StopProcessing != nil {
		update.StopProcessing = req.StopProcessing
	}

	return update
}

// mapRuleUpdateRequestToAPI converts RuleUpdateRequest to client.RuleUpdate
func mapRuleUpdateRequestToAPI(req *RuleUpdateRequest) client.RuleUpdate {
	update := client.RuleUpdate{
		Title:          req.Title,
		Description:    req.Description,
		RuleGroupId:    req.RuleGroupId,
		Active:         req.Active,
		Strict:         req.Strict,
		StopProcessing: req.StopProcessing,
	}

	// Map trigger type
	if req.Trigger != nil {
		triggerType := client.RuleTriggerType(*req.Trigger)
		update.Trigger = &triggerType
	}

	// Map triggers if provided
	if len(req.Triggers) > 0 {
		triggers := make([]client.RuleTriggerUpdate, len(req.Triggers))
		for i, trigger := range req.Triggers {
			triggers[i] = mapRuleTriggerRequestToUpdate(trigger)
		}
		update.Triggers = &triggers
	}

	// Map actions if provided
	if len(req.Actions) > 0 {
		actions := make([]client.RuleActionUpdate, len(req.Actions))
		for i, action := range req.Actions {
			actions[i] = mapRuleActionRequestToUpdate(action)
		}
		update.Actions = &actions
	}

	return update
}

// mapRuleGroupStoreRequestToAPI converts RuleGroupStoreRequest to client.RuleGroupStore
func mapRuleGroupStoreRequestToAPI(req *RuleGroupStoreRequest) client.RuleGroupStore {
	return client.RuleGroupStore{
		Title:       req.Title,
		Description: req.Description,
		Active:      req.Active,
	}
}

// mapRuleGroupUpdateRequestToAPI converts RuleGroupUpdateRequest to client.RuleGroupUpdate
func mapRuleGroupUpdateRequestToAPI(req *RuleGroupUpdateRequest) client.RuleGroupUpdate {
	return client.RuleGroupUpdate{
		Title:       req.Title,
		Description: req.Description,
		Active:      req.Active,
	}
}
