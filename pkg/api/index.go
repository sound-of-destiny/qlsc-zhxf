package api

import (
	"fmt"
	"strings"

	"github.com/sound-of-destiny/qlsc_zhxf/pkg/api/dtos"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/bus"
	m "github.com/sound-of-destiny/qlsc_zhxf/pkg/models"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/plugins"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/setting"
)

const (
	// Themes
	lightName = "light"
	darkName  = "dark"
)

func setIndexViewData(c *m.ReqContext) (*dtos.IndexViewData, error) {
	settings, err := getFrontendSettingsMap(c)
	if err != nil {
		return nil, err
	}

	prefsQuery := m.GetPreferencesWithDefaultsQuery{OrgId: c.OrgId, UserId: c.UserId}
	if err := bus.Dispatch(&prefsQuery); err != nil {
		return nil, err
	}
	prefs := prefsQuery.Result

	// Read locale from acccept-language
	acceptLang := c.Req.Header.Get("Accept-Language")
	locale := "en-US"

	if len(acceptLang) > 0 {
		parts := strings.Split(acceptLang, ",")
		locale = parts[0]
	}

	appURL := setting.AppUrl
	appSubURL := setting.AppSubUrl

	// special case when doing localhost call from phantomjs
	if c.IsRenderCall {
		appURL = fmt.Sprintf("%s://localhost:%s", setting.Protocol, setting.HttpPort)
		appSubURL = ""
		settings["appSubUrl"] = ""
	}

	hasEditPermissionInFoldersQuery := m.HasEditPermissionInFoldersQuery{SignedInUser: c.SignedInUser}
	if err := bus.Dispatch(&hasEditPermissionInFoldersQuery); err != nil {
		return nil, err
	}

	var data = dtos.IndexViewData{
		User: &dtos.CurrentUser{
			Id:                         c.UserId,
			IsSignedIn:                 c.IsSignedIn,
			Login:                      c.Login,
			Email:                      c.Email,
			Name:                       c.Name,
			OrgCount:                   c.OrgCount,
			OrgId:                      c.OrgId,
			OrgName:                    c.OrgName,
			OrgRole:                    c.OrgRole,
			GravatarUrl:                dtos.GetGravatarUrl(c.Email),
			IsGrafanaAdmin:             c.IsGrafanaAdmin,
			LightTheme:                 prefs.Theme == lightName,
			Timezone:                   prefs.Timezone,
			Locale:                     locale,
			HelpFlags1:                 c.HelpFlags1,
			HasEditPermissionInFolders: hasEditPermissionInFoldersQuery.Result,
		},
		Settings:                settings,
		Theme:                   prefs.Theme,
		AppUrl:                  appURL,
		AppSubUrl:               appSubURL,
		GoogleAnalyticsId:       setting.GoogleAnalyticsId,
		GoogleTagManagerId:      setting.GoogleTagManagerId,
		BuildVersion:            setting.BuildVersion,
		BuildCommit:             setting.BuildCommit,
		NewGrafanaVersion:       plugins.GrafanaLatestVersion,
		NewGrafanaVersionExists: plugins.GrafanaHasUpdate,
		AppName:                 setting.ApplicationName,
	}

	if setting.DisableGravatar {
		data.User.GravatarUrl = setting.AppSubUrl + "/public/img/user_profile.png"
	}

	if len(data.User.Name) == 0 {
		data.User.Name = data.User.Login
	}

	themeURLParam := c.Query("theme")
	if themeURLParam == lightName {
		data.User.LightTheme = true
		data.Theme = lightName
	} else if themeURLParam == darkName {
		data.User.LightTheme = false
		data.Theme = darkName
	}

	if hasEditPermissionInFoldersQuery.Result {
		children := []*dtos.NavLink{
			{Text: "Dashboard", Icon: "gicon gicon-dashboard-new", Url: setting.AppSubUrl + "/dashboard/new"},
		}

		if c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR {
			children = append(children, &dtos.NavLink{Text: "Folder", SubTitle: "Create a new folder to organize your dashboards", Id: "folder", Icon: "gicon gicon-folder-new", Url: setting.AppSubUrl + "/dashboards/folder/new"})
		}

		children = append(children, &dtos.NavLink{Text: "Import", SubTitle: "Import dashboard from file or Grafana.com", Id: "import", Icon: "gicon gicon-dashboard-import", Url: setting.AppSubUrl + "/dashboard/import"})

		data.NavTree = append(data.NavTree, &dtos.NavLink{
			Text:     "Create",
			Id:       "create",
			Icon:     "fa fa-fw fa-plus",
			Url:      setting.AppSubUrl + "/dashboard/new",
			Children: children,
		})
	}

	//原dashboard
	/*dashboardChildNavs := []*dtos.NavLink{
		{Text: "Home", Id: "home", Url: setting.AppSubUrl + "/", Icon: "gicon gicon-home", HideFromTabs: true},
		{Text: "Divider", Divider: true, Id: "divider", HideFromTabs: true},
		{Text: "Manage", Id: "manage-dashboards", Url: setting.AppSubUrl + "/dashboards", Icon: "gicon gicon-manage"},
		{Text: "Playlists", Id: "playlists", Url: setting.AppSubUrl + "/playlists", Icon: "gicon gicon-playlists"},
		{Text: "Snapshots", Id: "snapshots", Url: setting.AppSubUrl + "/dashboard/snapshots", Icon: "gicon gicon-snapshots"},
	}

	data.NavTree = append(data.NavTree, &dtos.NavLink{
		Text:     "主页",
		Id:       "dashboards",
		SubTitle: "Manage dashboards & folders",
		Icon:     "gicon gicon-dashboard",
		Url:      setting.AppSubUrl + "/",
		Children: dashboardChildNavs,
	})*/

	//主页
	data.NavTree = append(data.NavTree, &dtos.NavLink{
		Text:     "主页",
		Id:       "home",
		SubTitle: "Show main views",
		Icon:     "gicon gicon-home",
		Url:      setting.AppSubUrl + "/",
		Children: make([]*dtos.NavLink, 0),
	})

	//dashboard列表
	data.NavTree = append(data.NavTree, &dtos.NavLink{
		Text:     "监控列表",
		Id:       "manage-dashboards",
		SubTitle: "List of dashboard",
		Icon:     "gicon gicon-dashboard",
		Url:      setting.AppSubUrl + "/dashboards",
		Children: make([]*dtos.NavLink, 0),
	})

	//探索页面
	if setting.ExploreEnabled && (c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR) {
		data.NavTree = append(data.NavTree, &dtos.NavLink{
			Text:     "Explore",
			Id:       "explore",
			SubTitle: "Explore your data",
			Icon:     "fa fa-rocket",
			Url:      setting.AppSubUrl + "/explore",
			Children: []*dtos.NavLink{
				{Text: "New tab", Icon: "gicon gicon-dashboard-new", Url: setting.AppSubUrl + "/explore"},
			},
		})
	}

	//报警模块
	if setting.AlertingEnabled && (c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR) {
		alertChildNavs := []*dtos.NavLink{
			{Text: "报警规则", Id: "alert-list", Url: setting.AppSubUrl + "/alerting/list", Icon: "gicon gicon-alert-rules"},
			{Text: "提示方式", Id: "channels", Url: setting.AppSubUrl + "/alerting/notifications", Icon: "gicon gicon-alert-notification-channel"},
		}

		data.NavTree = append(data.NavTree, &dtos.NavLink{
			Text:     "报警",
			SubTitle: "Alert rules & notifications",
			Id:       "alerting",
			Icon:     "gicon gicon-alert",
			Url:      setting.AppSubUrl + "/alerting/list",
			Children: alertChildNavs,
		})
	}

	/*enabledPlugins, err := plugins.GetEnabledPlugins(c.OrgId)
	if err != nil {
		return nil, err
	}

	for _, plugin := range enabledPlugins.Apps {
		if plugin.Pinned {
			appLink := &dtos.NavLink{
				Text: plugin.Name,
				Id:   "plugin-page-" + plugin.Id,
				Url:  plugin.DefaultNavUrl,
				Img:  plugin.Info.Logos.Small,
			}

			for _, include := range plugin.Includes {
				if !c.HasUserRole(include.Role) {
					continue
				}

				if include.Type == "page" && include.AddToNav {
					link := &dtos.NavLink{
						Url:  setting.AppSubUrl + "/plugins/" + plugin.Id + "/page/" + include.Slug,
						Text: include.Name,
					}
					appLink.Children = append(appLink.Children, link)
				}

				if include.Type == "dashboard" && include.AddToNav {
					link := &dtos.NavLink{
						Url:  setting.AppSubUrl + "/dashboard/db/" + include.Slug,
						Text: include.Name,
					}
					appLink.Children = append(appLink.Children, link)
				}
			}

			if len(appLink.Children) > 0 && c.OrgRole == m.ROLE_ADMIN {
				appLink.Children = append(appLink.Children, &dtos.NavLink{Divider: true})
				appLink.Children = append(appLink.Children, &dtos.NavLink{Text: "Plugin Config", Icon: "gicon gicon-cog", Url: setting.AppSubUrl + "/plugins/" + plugin.Id + "/edit"})
			}

			if len(appLink.Children) > 0 {
				data.NavTree = append(data.NavTree, appLink)
			}
		}
	}*/

	//设置中心
	if c.IsGrafanaAdmin || c.OrgRole == m.ROLE_EDITOR {
		cfgNode := &dtos.NavLink{
			Id:       "cfg",
			Text:     "设置",
			SubTitle: "Organization: " + c.OrgName,
			Icon:     "gicon gicon-cog",
			Url:      setting.AppSubUrl + "/org/users",
			Children: []*dtos.NavLink{
				{
					Text:        "用户",
					Id:          "users",
					Description: "Manage org members",
					Icon:        "gicon gicon-user",
					Url:         setting.AppSubUrl + "/org/users",
				},
				/*{
					Text:        "数据源",
					Icon:        "gicon gicon-datasources",
					Description: "Add and configure data sources",
					Id:          "datasources",
					Url:         setting.AppSubUrl + "/datasources",
				},*/
			},
		}

		if c.IsGrafanaAdmin {
			cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
				Divider: true, HideFromTabs: true, Id: "admin-divider", Text: "Text",
			})

			cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
				Text:         "管理员",
				HideFromTabs: true,
				SubTitle:     "Manage all users",
				Id:           "admin",
				Icon:         "gicon gicon-shield",
				Url:          setting.AppSubUrl + "/admin/users",
				Children: []*dtos.NavLink{
					{Text: "用户", Id: "global-users", Url: setting.AppSubUrl + "/admin/users", Icon: "gicon gicon-user"},
					{Text: "数据源", Id: "datasources", Url: setting.AppSubUrl + "/datasources", Icon: "gicon gicon-datasources"},
					{Text: "设置", Id: "server-settings", Url: setting.AppSubUrl + "/admin/settings", Icon: "gicon gicon-preferences"},
					{Text: "状态", Id: "server-stats", Url: setting.AppSubUrl + "/admin/stats", Icon: "fa fa-fw fa-bar-chart"},
					{Text: "设计", Id: "styleguide", Url: setting.AppSubUrl + "/styleguide", Icon: "fa fa-fw fa-eyedropper"},
				},
			})
		}

		data.NavTree = append(data.NavTree, cfgNode)
	}

	//个人中心
	if c.IsSignedIn {
		// Only set login if it's different from the name
		var login string
		if c.SignedInUser.Login != c.SignedInUser.NameOrFallback() {
			login = c.SignedInUser.Login
		}
		profileNode := &dtos.NavLink{
			Text:         c.SignedInUser.NameOrFallback(),
			SubTitle:     login,
			Id:           "profile",
			Img:          data.User.GravatarUrl,
			Url:          setting.AppSubUrl + "/profile",
			HideFromMenu: true,
			Children: []*dtos.NavLink{
				{Text: "个人中心", Id: "profile-settings", Url: setting.AppSubUrl + "/profile", Icon: "gicon gicon-preferences"},
				{Text: "修改密码", Id: "change-password", Url: setting.AppSubUrl + "/profile/password", Icon: "fa fa-fw fa-lock", HideFromMenu: true},
			},
		}

		if !setting.DisableSignoutMenu {
			// add sign out first
			profileNode.Children = append(profileNode.Children, &dtos.NavLink{
				Text: "退出", Id: "sign-out", Url: setting.AppSubUrl + "/logout", Icon: "fa fa-fw fa-sign-out", Target: "_self",
			})
		}

		data.NavTree = append(data.NavTree, profileNode)
	}

	//帮助
	/*data.NavTree = append(data.NavTree, &dtos.NavLink{
		Text:         "Help",
		SubTitle:     fmt.Sprintf(`%s v%s (%s)`, setting.ApplicationName, setting.BuildVersion, setting.BuildCommit),
		Id:           "help",
		Url:          "#",
		Icon:         "gicon gicon-question",
		HideFromMenu: true,
		Children: []*dtos.NavLink{
			{Text: "Keyboard shortcuts", Url: "/shortcuts", Icon: "fa fa-fw fa-keyboard-o", Target: "_self"},
			{Text: "Community site", Url: "http://community.grafana.com", Icon: "fa fa-fw fa-comment", Target: "_blank"},
			{Text: "Documentation", Url: "http://docs.grafana.org", Icon: "fa fa-fw fa-file", Target: "_blank"},
		},
	})*/

	return &data, nil
}

func Index(c *m.ReqContext) {
	data, err := setIndexViewData(c)
	if err != nil {
		c.Handle(500, "Failed to get settings", err)
		return
	}
	c.HTML(200, "index", data)
}

func NotFoundHandler(c *m.ReqContext) {
	if c.IsApiRequest() {
		c.JsonApiErr(404, "Not found", nil)
		return
	}

	data, err := setIndexViewData(c)
	if err != nil {
		c.Handle(500, "Failed to get settings", err)
		return
	}

	c.HTML(404, "index", data)
}
