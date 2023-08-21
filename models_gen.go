// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package therapy

type TherapistFilters struct {
	Title                 *string `json:"title,omitempty"`
	AcceptingAppointments *bool   `json:"accepting_appointments,omitempty"`
	Credentials           *string `json:"credentials,omitempty"`
	Verified              *string `json:"verified,omitempty"`
	Statement             *string `json:"statement,omitempty"`
	Phone                 *string `json:"phone,omitempty"`
	Location              *string `json:"location,omitempty"`
	Link                  *string `json:"link,omitempty"`
	Limit                 *int    `json:"limit,omitempty"`
	Offset                *int    `json:"offset,omitempty"`
}