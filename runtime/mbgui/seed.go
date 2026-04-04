package mbgui

import "moonbasic/vm/value"

// SeedGUIGlobals installs raygui control/property IDs and related enums for GUI.SETSTYLE / GUI.SETCOLOR / GUI.GETTEXTBOUNDS.
// Full naming guide and every GUI.* signature: docs/reference/GUI.md.
// Numeric values match github.com/gen2brain/raylib-go/raygui (raygui 4.x layout).
// Several GPROP_* names share the same integer for different control kinds; always pair with GCTL_*.
func SeedGUIGlobals(globals map[string]value.Value) {
	// --- Control IDs (raygui.ControlID) ---
	globals["GCTL_DEFAULT"] = value.FromInt(0)
	globals["GCTL_LABEL"] = value.FromInt(1)
	globals["GCTL_BUTTON"] = value.FromInt(2)
	globals["GCTL_TOGGLE"] = value.FromInt(3)
	globals["GCTL_SLIDER"] = value.FromInt(4)
	globals["GCTL_PROGRESSBAR"] = value.FromInt(5)
	globals["GCTL_CHECKBOX"] = value.FromInt(6)
	globals["GCTL_COMBOBOX"] = value.FromInt(7)
	globals["GCTL_DROPDOWNBOX"] = value.FromInt(8)
	globals["GCTL_TEXTBOX"] = value.FromInt(9)
	globals["GCTL_VALUEBOX"] = value.FromInt(10)
	globals["GCTL_CONTROL11"] = value.FromInt(11)
	globals["GCTL_LISTVIEW"] = value.FromInt(12)
	globals["GCTL_COLORPICKER"] = value.FromInt(13)
	globals["GCTL_SCROLLBAR"] = value.FromInt(14)
	globals["GCTL_STATUSBAR"] = value.FromInt(15)

	// --- Base property IDs (colors, sizes) — use with matching control ---
	globals["GPROP_BORDER_COLOR_NORMAL"] = value.FromInt(0)
	globals["GPROP_BASE_COLOR_NORMAL"] = value.FromInt(1)
	globals["GPROP_TEXT_COLOR_NORMAL"] = value.FromInt(2)
	globals["GPROP_BORDER_COLOR_FOCUSED"] = value.FromInt(3)
	globals["GPROP_BASE_COLOR_FOCUSED"] = value.FromInt(4)
	globals["GPROP_TEXT_COLOR_FOCUSED"] = value.FromInt(5)
	globals["GPROP_BORDER_COLOR_PRESSED"] = value.FromInt(6)
	globals["GPROP_BASE_COLOR_PRESSED"] = value.FromInt(7)
	globals["GPROP_TEXT_COLOR_PRESSED"] = value.FromInt(8)
	globals["GPROP_BORDER_COLOR_DISABLED"] = value.FromInt(9)
	globals["GPROP_BASE_COLOR_DISABLED"] = value.FromInt(10)
	globals["GPROP_TEXT_COLOR_DISABLED"] = value.FromInt(11)
	globals["GPROP_BORDER_WIDTH"] = value.FromInt(12)
	globals["GPROP_TEXT_PADDING"] = value.FromInt(13)
	globals["GPROP_TEXT_ALIGNMENT"] = value.FromInt(14)

	// --- DEFAULT (global) extended ---
	globals["GPROP_TEXT_SIZE"] = value.FromInt(16)
	globals["GPROP_TEXT_SPACING"] = value.FromInt(17)
	globals["GPROP_LINE_COLOR"] = value.FromInt(18)
	globals["GPROP_BACKGROUND_COLOR"] = value.FromInt(19)
	globals["GPROP_TEXT_LINE_SPACING"] = value.FromInt(20)
	globals["GPROP_TEXT_ALIGNMENT_VERTICAL"] = value.FromInt(21)
	globals["GPROP_TEXT_WRAP_MODE"] = value.FromInt(22)

	// --- Per-control extended (values often overlap; control ID selects meaning) ---
	globals["GPROP_TOGGLE_GROUP_PADDING"] = value.FromInt(16)
	globals["GPROP_SLIDER_WIDTH"] = value.FromInt(16)
	globals["GPROP_SLIDER_PADDING"] = value.FromInt(17)
	globals["GPROP_PROGRESS_PADDING"] = value.FromInt(16)
	globals["GPROP_SCROLLBAR_ARROWS_SIZE"] = value.FromInt(16)
	globals["GPROP_SCROLLBAR_ARROWS_VISIBLE"] = value.FromInt(17)
	globals["GPROP_SCROLLBAR_SLIDER_PADDING"] = value.FromInt(18)
	globals["GPROP_SCROLLBAR_SLIDER_SIZE"] = value.FromInt(19)
	globals["GPROP_SCROLLBAR_SCROLL_PADDING"] = value.FromInt(20)
	globals["GPROP_SCROLLBAR_SCROLL_SPEED"] = value.FromInt(21)
	globals["GPROP_CHECK_PADDING"] = value.FromInt(16)
	globals["GPROP_COMBO_BUTTON_WIDTH"] = value.FromInt(16)
	globals["GPROP_COMBO_BUTTON_SPACING"] = value.FromInt(17)
	globals["GPROP_DROPDOWN_ARROW_PADDING"] = value.FromInt(16)
	globals["GPROP_DROPDOWN_ITEMS_SPACING"] = value.FromInt(17)
	globals["GPROP_DROPDOWN_ARROW_HIDDEN"] = value.FromInt(18)
	globals["GPROP_DROPDOWN_ROLL_UP"] = value.FromInt(19)
	globals["GPROP_TEXT_READONLY"] = value.FromInt(16)
	globals["GPROP_SPINNER_BUTTON_WIDTH"] = value.FromInt(16)
	globals["GPROP_SPINNER_BUTTON_SPACING"] = value.FromInt(17)
	globals["GPROP_LIST_ITEMS_HEIGHT"] = value.FromInt(16)
	globals["GPROP_LIST_ITEMS_SPACING"] = value.FromInt(17)
	globals["GPROP_LIST_SCROLLBAR_WIDTH"] = value.FromInt(18)
	globals["GPROP_LIST_SCROLLBAR_SIDE"] = value.FromInt(19)
	globals["GPROP_LIST_ITEMS_BORDER_NORMAL"] = value.FromInt(20)
	globals["GPROP_LIST_ITEMS_BORDER_WIDTH"] = value.FromInt(21)
	globals["GPROP_COLOR_SELECTOR_SIZE"] = value.FromInt(16)
	globals["GPROP_COLOR_HUEBAR_WIDTH"] = value.FromInt(17)
	globals["GPROP_COLOR_HUEBAR_PADDING"] = value.FromInt(18)
	globals["GPROP_COLOR_HUEBAR_SELECTOR_HEIGHT"] = value.FromInt(19)
	globals["GPROP_COLOR_HUEBAR_SELECTOR_OVERFLOW"] = value.FromInt(20)

	// --- Gui state (GUI.SETSTATE / focus drawing) ---
	globals["GUI_STATE_NORMAL"] = value.FromInt(0)
	globals["GUI_STATE_FOCUSED"] = value.FromInt(1)
	globals["GUI_STATE_PRESSED"] = value.FromInt(2)
	globals["GUI_STATE_DISABLED"] = value.FromInt(3)

	// --- Text alignment (horizontal / vertical / wrap) ---
	globals["GUI_TEXT_ALIGN_LEFT"] = value.FromInt(0)
	globals["GUI_TEXT_ALIGN_CENTER"] = value.FromInt(1)
	globals["GUI_TEXT_ALIGN_RIGHT"] = value.FromInt(2)
	globals["GUI_TEXT_ALIGN_TOP"] = value.FromInt(0)
	globals["GUI_TEXT_ALIGN_MIDDLE"] = value.FromInt(1)
	globals["GUI_TEXT_ALIGN_BOTTOM"] = value.FromInt(2)
	globals["GUI_TEXT_WRAP_NONE"] = value.FromInt(0)
	globals["GUI_TEXT_WRAP_CHAR"] = value.FromInt(1)
	globals["GUI_TEXT_WRAP_WORD"] = value.FromInt(2)

	// --- Scrollbar side ---
	globals["GUI_SCROLLBAR_LEFT"] = value.FromInt(0)
	globals["GUI_SCROLLBAR_RIGHT"] = value.FromInt(1)
}
