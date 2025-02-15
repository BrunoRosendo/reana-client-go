// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GetWorkflowLogsReader is a Reader for the GetWorkflowLogs structure.
type GetWorkflowLogsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetWorkflowLogsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetWorkflowLogsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetWorkflowLogsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewGetWorkflowLogsForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetWorkflowLogsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetWorkflowLogsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetWorkflowLogsOK creates a GetWorkflowLogsOK with default headers values
func NewGetWorkflowLogsOK() *GetWorkflowLogsOK {
	return &GetWorkflowLogsOK{}
}

/* GetWorkflowLogsOK describes a response with status code 200, with default header values.

Request succeeded. Info about a workflow, including the status is returned.
*/
type GetWorkflowLogsOK struct {
	Payload *GetWorkflowLogsOKBody
}

func (o *GetWorkflowLogsOK) Error() string {
	return fmt.Sprintf("[GET /api/workflows/{workflow_id_or_name}/logs][%d] getWorkflowLogsOK  %+v", 200, o.Payload)
}
func (o *GetWorkflowLogsOK) GetPayload() *GetWorkflowLogsOKBody {
	return o.Payload
}

func (o *GetWorkflowLogsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetWorkflowLogsOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkflowLogsBadRequest creates a GetWorkflowLogsBadRequest with default headers values
func NewGetWorkflowLogsBadRequest() *GetWorkflowLogsBadRequest {
	return &GetWorkflowLogsBadRequest{}
}

/* GetWorkflowLogsBadRequest describes a response with status code 400, with default header values.

Request failed. The incoming data specification seems malformed.
*/
type GetWorkflowLogsBadRequest struct {
	Payload *GetWorkflowLogsBadRequestBody
}

func (o *GetWorkflowLogsBadRequest) Error() string {
	return fmt.Sprintf("[GET /api/workflows/{workflow_id_or_name}/logs][%d] getWorkflowLogsBadRequest  %+v", 400, o.Payload)
}
func (o *GetWorkflowLogsBadRequest) GetPayload() *GetWorkflowLogsBadRequestBody {
	return o.Payload
}

func (o *GetWorkflowLogsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetWorkflowLogsBadRequestBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkflowLogsForbidden creates a GetWorkflowLogsForbidden with default headers values
func NewGetWorkflowLogsForbidden() *GetWorkflowLogsForbidden {
	return &GetWorkflowLogsForbidden{}
}

/* GetWorkflowLogsForbidden describes a response with status code 403, with default header values.

Request failed. User is not allowed to access workflow.
*/
type GetWorkflowLogsForbidden struct {
	Payload *GetWorkflowLogsForbiddenBody
}

func (o *GetWorkflowLogsForbidden) Error() string {
	return fmt.Sprintf("[GET /api/workflows/{workflow_id_or_name}/logs][%d] getWorkflowLogsForbidden  %+v", 403, o.Payload)
}
func (o *GetWorkflowLogsForbidden) GetPayload() *GetWorkflowLogsForbiddenBody {
	return o.Payload
}

func (o *GetWorkflowLogsForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetWorkflowLogsForbiddenBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkflowLogsNotFound creates a GetWorkflowLogsNotFound with default headers values
func NewGetWorkflowLogsNotFound() *GetWorkflowLogsNotFound {
	return &GetWorkflowLogsNotFound{}
}

/* GetWorkflowLogsNotFound describes a response with status code 404, with default header values.

Request failed. User does not exist.
*/
type GetWorkflowLogsNotFound struct {
	Payload *GetWorkflowLogsNotFoundBody
}

func (o *GetWorkflowLogsNotFound) Error() string {
	return fmt.Sprintf("[GET /api/workflows/{workflow_id_or_name}/logs][%d] getWorkflowLogsNotFound  %+v", 404, o.Payload)
}
func (o *GetWorkflowLogsNotFound) GetPayload() *GetWorkflowLogsNotFoundBody {
	return o.Payload
}

func (o *GetWorkflowLogsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetWorkflowLogsNotFoundBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWorkflowLogsInternalServerError creates a GetWorkflowLogsInternalServerError with default headers values
func NewGetWorkflowLogsInternalServerError() *GetWorkflowLogsInternalServerError {
	return &GetWorkflowLogsInternalServerError{}
}

/* GetWorkflowLogsInternalServerError describes a response with status code 500, with default header values.

Request failed. Internal controller error.
*/
type GetWorkflowLogsInternalServerError struct {
	Payload *GetWorkflowLogsInternalServerErrorBody
}

func (o *GetWorkflowLogsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/workflows/{workflow_id_or_name}/logs][%d] getWorkflowLogsInternalServerError  %+v", 500, o.Payload)
}
func (o *GetWorkflowLogsInternalServerError) GetPayload() *GetWorkflowLogsInternalServerErrorBody {
	return o.Payload
}

func (o *GetWorkflowLogsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetWorkflowLogsInternalServerErrorBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetWorkflowLogsBadRequestBody get workflow logs bad request body
swagger:model GetWorkflowLogsBadRequestBody
*/
type GetWorkflowLogsBadRequestBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this get workflow logs bad request body
func (o *GetWorkflowLogsBadRequestBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get workflow logs bad request body based on context it is used
func (o *GetWorkflowLogsBadRequestBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetWorkflowLogsBadRequestBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetWorkflowLogsBadRequestBody) UnmarshalBinary(b []byte) error {
	var res GetWorkflowLogsBadRequestBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*GetWorkflowLogsForbiddenBody get workflow logs forbidden body
swagger:model GetWorkflowLogsForbiddenBody
*/
type GetWorkflowLogsForbiddenBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this get workflow logs forbidden body
func (o *GetWorkflowLogsForbiddenBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get workflow logs forbidden body based on context it is used
func (o *GetWorkflowLogsForbiddenBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetWorkflowLogsForbiddenBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetWorkflowLogsForbiddenBody) UnmarshalBinary(b []byte) error {
	var res GetWorkflowLogsForbiddenBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*GetWorkflowLogsInternalServerErrorBody get workflow logs internal server error body
swagger:model GetWorkflowLogsInternalServerErrorBody
*/
type GetWorkflowLogsInternalServerErrorBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this get workflow logs internal server error body
func (o *GetWorkflowLogsInternalServerErrorBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get workflow logs internal server error body based on context it is used
func (o *GetWorkflowLogsInternalServerErrorBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetWorkflowLogsInternalServerErrorBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetWorkflowLogsInternalServerErrorBody) UnmarshalBinary(b []byte) error {
	var res GetWorkflowLogsInternalServerErrorBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*GetWorkflowLogsNotFoundBody get workflow logs not found body
swagger:model GetWorkflowLogsNotFoundBody
*/
type GetWorkflowLogsNotFoundBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this get workflow logs not found body
func (o *GetWorkflowLogsNotFoundBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get workflow logs not found body based on context it is used
func (o *GetWorkflowLogsNotFoundBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetWorkflowLogsNotFoundBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetWorkflowLogsNotFoundBody) UnmarshalBinary(b []byte) error {
	var res GetWorkflowLogsNotFoundBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*GetWorkflowLogsOKBody get workflow logs o k body
swagger:model GetWorkflowLogsOKBody
*/
type GetWorkflowLogsOKBody struct {

	// logs
	Logs string `json:"logs,omitempty"`

	// user
	User string `json:"user,omitempty"`

	// workflow id
	WorkflowID string `json:"workflow_id,omitempty"`

	// workflow name
	WorkflowName string `json:"workflow_name,omitempty"`
}

// Validate validates this get workflow logs o k body
func (o *GetWorkflowLogsOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get workflow logs o k body based on context it is used
func (o *GetWorkflowLogsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetWorkflowLogsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetWorkflowLogsOKBody) UnmarshalBinary(b []byte) error {
	var res GetWorkflowLogsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
