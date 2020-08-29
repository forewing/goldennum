package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// TemplateBaseURL is the base URL template name
	TemplateBaseURL = "base_url"

	templateHeader     = "header.html"
	templateFooter     = "footer.html"
	templateNavbar     = "navbar.html"
	templateIndex      = "index.html"
	templateUserPanel  = "user_panel.html"
	templateUserModals = "user_modals.html"
	templateAdmin      = "admin.html"

	componentDashboard   = "dashboard.html"
	componentRoomControl = "room_control.html"

	versionInfo = "version.html"
)

var (
	// Templates is the template names
	Templates []string = []string{
		templateHeader,
		templateFooter,
		templateNavbar,
		templateIndex,
		templateUserPanel,
		templateUserModals,
		templateAdmin,

		componentDashboard,
		componentRoomControl,

		versionInfo,
	}
)

// PageIndex render index
func PageIndex(c *gin.Context) {
	c.HTML(http.StatusOK, templateIndex, gin.H{})
}

// AdminIndex return admin page
func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, templateAdmin, gin.H{})
}
