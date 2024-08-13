package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/slack-go/slack"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

type SlackHandler struct {
	SigningSecret string
	SlackBotToken string
	DynamicClient dynamic.Interface
}

func NewSlackHandler(signingSecret, slackBotToken string) (*SlackHandler, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "/path/to/your/kubeconfig")
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %v", err)
	}

	return &SlackHandler{
		SigningSecret: signingSecret,
		SlackBotToken: slackBotToken,
		DynamicClient: dynamicClient,
	}, nil
}

func (sh *SlackHandler) CommandHandler(w http.ResponseWriter, r *http.Request) {
	verifier, err := slack.NewSecretsVerifier(r.Header, sh.SigningSecret)
	if err != nil {
		http.Error(w, "Failed to create secret verifier", http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		http.Error(w, "Failed to parse slash command", http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		http.Error(w, "Failed to verify request", http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/modify-power-config":
		sh.handleModifyPowerConfig(w)
	default:
		http.Error(w, "Unknown command", http.StatusBadRequest)
	}
}

func (sh *SlackHandler) InteractionHandler(w http.ResponseWriter, r *http.Request) {
	verifier, err := slack.NewSecretsVerifier(r.Header, sh.SigningSecret)
	if err != nil {
		http.Error(w, "Failed to create secret verifier", http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(io.TeeReader(r.Body, &verifier))

	var payload slack.InteractionCallback
	err = json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		http.Error(w, "Failed to parse interaction payload", http.StatusBadRequest)
		return
	}

	if err = verifier.Ensure(); err != nil {
		http.Error(w, "Failed to verify request", http.StatusUnauthorized)
		return
	}

	switch payload.Type {
	case slack.InteractionTypeBlockActions:
		sh.handleBlockActions(w, &payload)
	case slack.InteractionTypeViewSubmission:
		sh.handleViewSubmission(w, &payload)
	default:
		http.Error(w, "Unknown interaction type", http.StatusBadRequest)
	}
}

func (sh *SlackHandler) handleBlockActions(w http.ResponseWriter, payload *slack.InteractionCallback) {
	for _, action := range payload.ActionCallback.BlockActions {
		switch action.ActionID {
		case "select_parameter":
			parameter := action.SelectedOption.Value
			sh.openParameterInputModal(w, payload.TriggerID, parameter)
		default:
			http.Error(w, "Unknown action", http.StatusBadRequest)
		}
	}
}

func (sh *SlackHandler) handleViewSubmission(w http.ResponseWriter, payload *slack.InteractionCallback) {
	switch payload.View.CallbackID {
	case "set_efficiency_level", "set_power_cap_percentage":
		sh.handleParameterUpdate(w, payload)
	default:
		http.Error(w, "Unknown view submission", http.StatusBadRequest)
	}
}

func (sh *SlackHandler) openParameterInputModal(w http.ResponseWriter, triggerID, parameter string) {
	modalView := slack.ModalViewRequest{
		Type: slack.ViewType("modal"),
		Title: &slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "Set " + parameter,
		},
		Submit: &slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "Submit",
		},
		Close: &slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "Cancel",
		},
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.InputBlock{
					BlockID: "new_value_block",
					Label: &slack.TextBlockObject{
						Type: slack.PlainTextType,
						Text: "New Value",
					},
					Element: slack.PlainTextInputBlockElement{
						Type:     slack.METPlainTextInput,
						ActionID: "new_value_action",
					},
				},
			},
		},
		CallbackID: "set_" + parameter,
	}

	api := slack.New(sh.SlackBotToken)
	_, err := api.OpenView(triggerID, modalView)
	if err != nil {
		http.Error(w, "Failed to open modal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (sh *SlackHandler) handleParameterUpdate(w http.ResponseWriter, payload *slack.InteractionCallback) {
	param := payload.View.CallbackID[4:] // Remove "set_" prefix
	value := payload.View.State.Values["new_value_block"]["new_value_action"].Value

	err := sh.UpdatePowerCappingConfig(param, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseMessage := fmt.Sprintf("Successfully updated %s to %s", param, value)

	api := slack.New(sh.SlackBotToken)
	_, _, err = api.PostMessage(
		payload.User.ID,
		slack.MsgOptionText(responseMessage, false),
	)
	if err != nil {
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response_action": "clear",
	})
}

func (sh *SlackHandler) handleModifyPowerConfig(w http.ResponseWriter) {
	blocks := []slack.Block{
		slack.NewSectionBlock(
			&slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: "Choose a parameter to modify:",
			},
			nil,
			nil,
		),
		slack.NewActionBlock(
			"parameter_selection",
			slack.NewOptionsSelectBlockElement(
				slack.OptTypeStatic,
				slack.NewTextBlockObject(slack.PlainTextType, "Select a parameter", false, false),
				"select_parameter",
				slack.NewOptionBlockObject(
					"efficiency_level",
					slack.NewTextBlockObject(slack.PlainTextType, "Efficiency Level", false, false),
					slack.NewTextBlockObject(slack.PlainTextType, "Set the efficiency level", false, false),
				),
				slack.NewOptionBlockObject(
					"power_cap_percentage",
					slack.NewTextBlockObject(slack.PlainTextType, "Power Cap Percentage", false, false),
					slack.NewTextBlockObject(slack.PlainTextType, "Set the power cap percentage", false, false),
				),
			),
		),
	}

	msg := slack.NewBlockMessage(blocks...)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func (sh *SlackHandler) UpdatePowerCappingConfig(param, value string) error {
	gvr := schema.GroupVersionResource{
		Group:    "climatik-project.io",
		Version:  "v1alpha1",
		Resource: "powercappingconfigs",
	}

	// Get the PowerCappingConfig
	pcc, err := sh.DynamicClient.Resource(gvr).Namespace("operator-powercapping-system").Get(context.Background(), "high-efficiency-for-stress-powercappingconfig", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get PowerCappingConfig: %v", err)
	}

	// Update the specified parameter
	spec := pcc.Object["spec"].(map[string]interface{})
	switch param {
	case "efficiency_level":
		spec["efficiencyLevel"] = value
	case "power_cap_percentage":
		percentage, _ := strconv.Atoi(value)
		powerCappingSpec := spec["powerCappingSpec"].(map[string]interface{})
		relativePowerCapSpec := powerCappingSpec["relativePowerCapInPercentageSpec"].(map[string]interface{})
		relativePowerCapSpec["powerCapPercentage"] = percentage
	default:
		return fmt.Errorf("unknown parameter: %s", param)
	}

	// Update the PowerCappingConfig in Kubernetes
	_, err = sh.DynamicClient.Resource(gvr).Namespace("default").Update(context.Background(), pcc, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update PowerCappingConfig: %v", err)
	}

	return nil
}
