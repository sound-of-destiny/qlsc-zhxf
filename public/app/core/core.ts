import './directives/dash_class';
import './directives/dropdown_typeahead';
import './directives/metric_segment';
import './directives/misc';
import './directives/ng_model_on_blur';
import './directives/tags';
import './directives/value_select_dropdown';
import './directives/rebuild_on_change';
import './directives/give_focus';
import './directives/diff-view';
import './jquery_extended';
import './partials';
import './utils/outline';
import './components/jsontree/jsontree';
import './components/code_editor/code_editor';
import './components/colorpicker/ColorPicker';
import './components/colorpicker/SeriesColorPicker';
import './components/colorpicker/spectrum_picker';
import './services/search_srv';
import './services/ng_react';

import 'app/core/controllers/all';
import 'app/core/services/all';
import './filters/filters';

import coreModule from './core_module';
import appEvents from './app_events';
import { arrayJoin } from './directives/array_join';
import { registerAngularDirectives } from './angular_wrappers';
import { Emitter } from './utils/emitter';
import colors from './utils/colors';
import { assignModelProperties } from './utils/model_utils';
import { profiler } from './profiler';
import { contextSrv } from './services/context_srv';
import { KeybindingSrv } from './services/keybindingSrv';
import { updateLegendValues } from './time_series2';
import TimeSeries from './time_series2';
import { NavModelSrv, NavModel } from './nav_model_srv';
import { grafanaAppDirective } from './components/grafana_app';
import { manageDashboardsDirective } from './components/manage_dashboards/manage_dashboards';
import { geminiScrollbar } from './components/scroll/scroll';
import { pageScrollbar } from './components/scroll/page_scroll';
import { gfPageDirective } from './components/gf_page';
import { JsonExplorer } from './components/json_explorer/json_explorer';
import { searchDirective } from './components/search/search';
import { layoutSelector } from './components/layout_selector/layout_selector';
import { sqlPartEditorDirective } from './components/sql_part/sql_part_editor';
import { searchResultsDirective } from './components/search/search_results';
import { navbarDirective } from './components/navbar/navbar';
import { dashboardSelector } from './components/dashboard_selector';
import { queryPartEditorDirective } from './components/query_part/query_part_editor';
import { formDropdownDirective } from './components/form_dropdown/form_dropdown';

/*import { liveSrv } from './live/live_srv';
import { infoPopover } from './components/info_popover';
import { switchDirective } from './components/switch';

import { helpModal } from './components/help/help';
import { orgSwitcher } from './components/org_switcher';
*/

export {
  coreModule,
  registerAngularDirectives,
  profiler,
  arrayJoin,
  grafanaAppDirective,
  navbarDirective,
  searchDirective,
  //liveSrv,
  layoutSelector,
  //switchDirective,
  //infoPopover,
  Emitter,
  appEvents,
  dashboardSelector,
  queryPartEditorDirective,
  sqlPartEditorDirective,
  colors,
  formDropdownDirective,
  assignModelProperties,
  contextSrv,
  KeybindingSrv,
  //helpModal,
  JsonExplorer,
  NavModelSrv,
  NavModel,
  geminiScrollbar,
  pageScrollbar,
  gfPageDirective,
  //orgSwitcher,
  manageDashboardsDirective,
  TimeSeries,
  updateLegendValues,
  searchResultsDirective,
};
