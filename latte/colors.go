package latte

// Named true-color palette for oat-latte.
//
// Colors are grouped by family and named after their role or source palette
// so themes and application code can reference semantic constants instead of
// raw hex literals. All values require a terminal with true-color support;
// fall back to the ANSI-16 variables in style.go for broad compatibility.
//
// Naming convention:
//
//	<Family><Shade>   e.g. SlateGray700, EmeraldGreen400
//
// Shade numbers follow the Tailwind CSS convention:
//
//	50   lightest tint
//	100 – 300  light
//	400 – 600  mid
//	700 – 900  dark
//	950  darkest shade
//
// Each palette also provides a semantic alias for the most common use-cases:
//
//	<Family>Light / <Family>Base / <Family>Dark

// ── Neutrals ─────────────────────────────────────────────────────────────────

var (
	// Slate — cool blue-gray scale
	Slate50  = RGB(248, 250, 252)
	Slate100 = RGB(241, 245, 249)
	Slate200 = RGB(226, 232, 240)
	Slate300 = RGB(203, 213, 225)
	Slate400 = RGB(148, 163, 184)
	Slate500 = RGB(100, 116, 139)
	Slate600 = RGB(71, 85, 105)
	Slate700 = RGB(51, 65, 85)
	Slate800 = RGB(30, 41, 59)
	Slate900 = RGB(15, 23, 42)
	Slate950 = RGB(2, 6, 23)

	// Zinc — neutral gray scale
	Zinc50  = RGB(250, 250, 250)
	Zinc100 = RGB(244, 244, 245)
	Zinc200 = RGB(228, 228, 231)
	Zinc300 = RGB(212, 212, 216)
	Zinc400 = RGB(161, 161, 170)
	Zinc500 = RGB(113, 113, 122)
	Zinc600 = RGB(82, 82, 91)
	Zinc700 = RGB(63, 63, 70)
	Zinc800 = RGB(39, 39, 42)
	Zinc900 = RGB(24, 24, 27)
	Zinc950 = RGB(9, 9, 11)

	// Stone — warm gray scale
	Stone50  = RGB(250, 250, 249)
	Stone100 = RGB(245, 245, 244)
	Stone200 = RGB(231, 229, 228)
	Stone300 = RGB(214, 211, 209)
	Stone400 = RGB(168, 162, 158)
	Stone500 = RGB(120, 113, 108)
	Stone600 = RGB(87, 83, 78)
	Stone700 = RGB(68, 64, 60)
	Stone800 = RGB(41, 37, 36)
	Stone900 = RGB(28, 25, 23)
	Stone950 = RGB(12, 10, 9)
)

// ── Blues ─────────────────────────────────────────────────────────────────────

var (
	// Sky — bright, airy sky blue
	Sky100 = RGB(224, 242, 254)
	Sky200 = RGB(186, 230, 253)
	Sky300 = RGB(125, 211, 252)
	Sky400 = RGB(56, 189, 248)
	Sky500 = RGB(14, 165, 233)
	Sky600 = RGB(2, 132, 199)
	Sky700 = RGB(3, 105, 161)
	Sky800 = RGB(7, 89, 133)
	Sky900 = RGB(12, 74, 110)

	// Blue — standard blue
	Blue100 = RGB(219, 234, 254)
	Blue200 = RGB(191, 219, 254)
	Blue300 = RGB(147, 197, 253)
	Blue400 = RGB(96, 165, 250)
	Blue500 = RGB(59, 130, 246)
	Blue600 = RGB(37, 99, 235)
	Blue700 = RGB(29, 78, 216)
	Blue800 = RGB(30, 64, 175)
	Blue900 = RGB(30, 58, 138)

	// Indigo — deep blue-violet
	Indigo100 = RGB(224, 231, 255)
	Indigo200 = RGB(199, 210, 254)
	Indigo300 = RGB(165, 180, 252)
	Indigo400 = RGB(129, 140, 248)
	Indigo500 = RGB(99, 102, 241)
	Indigo600 = RGB(79, 70, 229)
	Indigo700 = RGB(67, 56, 202)
	Indigo800 = RGB(55, 48, 163)
	Indigo900 = RGB(49, 46, 129)

	// Cornflower — the classic TUI "focus" blue
	Cornflower300 = RGB(147, 182, 255)
	Cornflower400 = RGB(100, 149, 237) // HTML cornflowerblue
	Cornflower500 = RGB(68, 114, 196)
	Cornflower600 = RGB(50, 90, 168)

	// Cyan — electric cyan / teal-leaning
	Cyan100 = RGB(207, 250, 254)
	Cyan200 = RGB(165, 243, 252)
	Cyan300 = RGB(103, 232, 249)
	Cyan400 = RGB(34, 211, 238)
	Cyan500 = RGB(6, 182, 212)
	Cyan600 = RGB(8, 145, 178)
	Cyan700 = RGB(14, 116, 144)
	Cyan800 = RGB(21, 94, 117)
	Cyan900 = RGB(22, 78, 99)

	// Teal — blue-green
	Teal100 = RGB(204, 251, 241)
	Teal200 = RGB(153, 246, 228)
	Teal300 = RGB(94, 234, 212)
	Teal400 = RGB(45, 212, 191)
	Teal500 = RGB(20, 184, 166)
	Teal600 = RGB(13, 148, 136)
	Teal700 = RGB(15, 118, 110)
	Teal800 = RGB(17, 94, 89)
	Teal900 = RGB(19, 78, 74)
)

// ── Greens ───────────────────────────────────────────────────────────────────

var (
	// Emerald — rich, vivid green
	Emerald100 = RGB(209, 250, 229)
	Emerald200 = RGB(167, 243, 208)
	Emerald300 = RGB(110, 231, 183)
	Emerald400 = RGB(52, 211, 153)
	Emerald500 = RGB(16, 185, 129)
	Emerald600 = RGB(5, 150, 105)
	Emerald700 = RGB(4, 120, 87)
	Emerald800 = RGB(6, 95, 70)
	Emerald900 = RGB(6, 78, 59)

	// Green — classic green
	Green100 = RGB(220, 252, 231)
	Green200 = RGB(187, 247, 208)
	Green300 = RGB(134, 239, 172)
	Green400 = RGB(74, 222, 128)
	Green500 = RGB(34, 197, 94)
	Green600 = RGB(22, 163, 74)
	Green700 = RGB(21, 128, 61)
	Green800 = RGB(22, 101, 52)
	Green900 = RGB(20, 83, 45)

	// Lime — yellow-green
	Lime300 = RGB(190, 242, 100)
	Lime400 = RGB(163, 230, 53)
	Lime500 = RGB(132, 204, 22)
	Lime600 = RGB(101, 163, 13)
	Lime700 = RGB(77, 124, 15)
)

// ── Yellows / Oranges ────────────────────────────────────────────────────────

var (
	// Yellow — warm yellow
	Yellow100 = RGB(254, 249, 195)
	Yellow200 = RGB(254, 240, 138)
	Yellow300 = RGB(253, 224, 71)
	Yellow400 = RGB(250, 204, 21)
	Yellow500 = RGB(234, 179, 8)
	Yellow600 = RGB(202, 138, 4)
	Yellow700 = RGB(161, 98, 7)
	Yellow800 = RGB(133, 77, 14)
	Yellow900 = RGB(113, 63, 18)

	// Amber — deep golden amber
	Amber100 = RGB(254, 243, 199)
	Amber200 = RGB(253, 230, 138)
	Amber300 = RGB(252, 211, 77)
	Amber400 = RGB(251, 191, 36)
	Amber500 = RGB(245, 158, 11)
	Amber600 = RGB(217, 119, 6)
	Amber700 = RGB(180, 83, 9)
	Amber800 = RGB(146, 64, 14)
	Amber900 = RGB(120, 53, 15)

	// Orange — vivid orange
	Orange100 = RGB(255, 237, 213)
	Orange200 = RGB(254, 215, 170)
	Orange300 = RGB(253, 186, 116)
	Orange400 = RGB(251, 146, 60)
	Orange500 = RGB(249, 115, 22)
	Orange600 = RGB(234, 88, 12)
	Orange700 = RGB(194, 65, 12)
	Orange800 = RGB(154, 52, 18)
	Orange900 = RGB(124, 45, 18)
)

// ── Reds / Pinks ─────────────────────────────────────────────────────────────

var (
	// Red — classic red
	Red100 = RGB(254, 226, 226)
	Red200 = RGB(254, 202, 202)
	Red300 = RGB(252, 165, 165)
	Red400 = RGB(248, 113, 113)
	Red500 = RGB(239, 68, 68)
	Red600 = RGB(220, 38, 38)
	Red700 = RGB(185, 28, 28)
	Red800 = RGB(153, 27, 27)
	Red900 = RGB(127, 29, 29)

	// Rose — warm pinkish red
	Rose100 = RGB(255, 228, 230)
	Rose200 = RGB(254, 205, 211)
	Rose300 = RGB(253, 164, 175)
	Rose400 = RGB(251, 113, 133)
	Rose500 = RGB(244, 63, 94)
	Rose600 = RGB(225, 29, 72)
	Rose700 = RGB(190, 18, 60)
	Rose800 = RGB(159, 18, 57)
	Rose900 = RGB(136, 19, 55)

	// Pink — bright pink
	Pink100 = RGB(252, 231, 243)
	Pink200 = RGB(251, 207, 232)
	Pink300 = RGB(249, 168, 212)
	Pink400 = RGB(244, 114, 182)
	Pink500 = RGB(236, 72, 153)
	Pink600 = RGB(219, 39, 119)
	Pink700 = RGB(190, 24, 93)
	Pink800 = RGB(157, 23, 77)
	Pink900 = RGB(131, 24, 67)
)

// ── Purples / Violets ────────────────────────────────────────────────────────

var (
	// Violet
	Violet100 = RGB(237, 233, 254)
	Violet200 = RGB(221, 214, 254)
	Violet300 = RGB(196, 181, 253)
	Violet400 = RGB(167, 139, 250)
	Violet500 = RGB(139, 92, 246)
	Violet600 = RGB(124, 58, 237)
	Violet700 = RGB(109, 40, 217)
	Violet800 = RGB(91, 33, 182)
	Violet900 = RGB(76, 29, 149)

	// Purple
	Purple100 = RGB(243, 232, 255)
	Purple200 = RGB(233, 213, 255)
	Purple300 = RGB(216, 180, 254)
	Purple400 = RGB(192, 132, 252)
	Purple500 = RGB(168, 85, 247)
	Purple600 = RGB(147, 51, 234)
	Purple700 = RGB(126, 34, 206)
	Purple800 = RGB(107, 33, 168)
	Purple900 = RGB(88, 28, 135)

	// Fuchsia — vivid magenta-purple
	Fuchsia100 = RGB(253, 244, 255)
	Fuchsia200 = RGB(245, 208, 254)
	Fuchsia300 = RGB(240, 171, 252)
	Fuchsia400 = RGB(232, 121, 249)
	Fuchsia500 = RGB(217, 70, 239)
	Fuchsia600 = RGB(192, 38, 211)
	Fuchsia700 = RGB(162, 28, 175)
	Fuchsia800 = RGB(134, 25, 143)
	Fuchsia900 = RGB(112, 26, 117)
)

// ── Design-system palettes ───────────────────────────────────────────────────
//
// The following groups are extracted from popular design systems / themes and
// are used directly by the built-in oat-latte themes.

// Dark theme palette — deep navy-black with blue-cyan accent
var (
	DarkBg         = RGB(15, 17, 23)    // #0f1117 — canvas / body background
	DarkBgElevated = RGB(26, 31, 46)    // #1a1f2e — dialog / elevated surface
	DarkBgScrim    = RGB(8, 11, 16)     // #080b10 — full-screen scrim
	DarkBorder     = RGB(46, 50, 71)    // #2e3247 — panel border (inactive)
	DarkMuted      = RGB(74, 80, 104)   // #4a5068 — muted text / inactive border
	DarkAccent     = RGB(124, 156, 255) // #7c9cff — primary accent (blue)
	DarkSuccess    = RGB(74, 222, 128)  // #4ade80 — green-400
	DarkWarning    = RGB(250, 204, 21)  // #facc15 — yellow-400
	DarkError      = RGB(248, 113, 113) // #f87171 — red-400
)

// Light theme palette — warm off-white with cobalt accent
var (
	LightBg         = RGB(248, 245, 240) // #f8f5f0 — canvas / body background
	LightBgElevated = RGB(238, 235, 229) // #eeebe5 — dialog / elevated surface
	LightBgSelected = RGB(219, 228, 253) // #dbe4fd — list row selection highlight (light cobalt tint)
	LightBgHeader   = RGB(232, 229, 224) // #e8e5e0 — header / footer surface (slightly darker than canvas)
	LightBgScrim    = RGB(180, 175, 168) // #b4afa8 — full-screen scrim (darkened for visible dimming)
	LightBorder     = RGB(200, 196, 190) // #c8c4be — panel border (inactive)
	LightMuted      = RGB(110, 110, 125) // #6e6e7d — muted / secondary text (darkened for WCAG AA contrast)
	LightText       = RGB(26, 26, 46)    // #1a1a2e — primary text
	LightAccent     = RGB(67, 97, 238)   // #4361ee — cobalt accent
	LightSuccess    = RGB(45, 106, 79)   // #2d6a4f — forest green
	LightWarning    = RGB(224, 123, 0)   // #e07b00 — amber
	LightError      = RGB(214, 40, 40)   // #d62828 — red
)

// Dracula palette — https://draculatheme.com
var (
	DraculaBg         = RGB(40, 42, 54)    // #282a36
	DraculaBgElevated = RGB(30, 31, 41)    // #1e1f29
	DraculaBgScrim    = RGB(26, 27, 36)    // #1a1b24
	DraculaFg         = RGB(248, 248, 242) // #f8f8f2
	DraculaComment    = RGB(98, 114, 164)  // #6272a4
	DraculaSelection  = RGB(68, 71, 90)    // #44475a
	DraculaPurple     = RGB(189, 147, 249) // #bd93f9
	DraculaCyan       = RGB(139, 233, 253) // #8be9fd
	DraculaGreen      = RGB(80, 250, 123)  // #50fa7b
	DraculaOrange     = RGB(255, 184, 108) // #ffb86c
	DraculaRed        = RGB(255, 85, 85)   // #ff5555
	DraculaYellow     = RGB(241, 250, 140) // #f1fa8c
	DraculaPink       = RGB(255, 121, 198) // #ff79c6
)

// Nord palette — https://www.nordtheme.com
var (
	// Polar Night (darkest → lightest)
	Nord0 = RGB(46, 52, 64)  // #2e3440
	Nord1 = RGB(59, 66, 82)  // #3b4252
	Nord2 = RGB(67, 76, 94)  // #434c5e
	Nord3 = RGB(76, 86, 106) // #4c566a

	// Snow Storm (lightest → darkest)
	Nord4 = RGB(216, 222, 233) // #d8dee9
	Nord5 = RGB(229, 233, 240) // #e5e9f0
	Nord6 = RGB(236, 239, 244) // #eceff4

	// Frost (blue accent family)
	Nord7  = RGB(143, 188, 187) // #8fbcbb — teal
	Nord8  = RGB(136, 192, 208) // #88c0d0 — light blue
	Nord9  = RGB(129, 161, 193) // #81a1c1 — blue
	Nord10 = RGB(94, 129, 172)  // #5e81ac — dark blue

	// Aurora (status / semantic)
	Nord11 = RGB(191, 97, 106)  // #bf616a — red
	Nord12 = RGB(208, 135, 112) // #d08770 — orange
	Nord13 = RGB(235, 203, 139) // #ebcb8b — yellow
	Nord14 = RGB(163, 190, 140) // #a3be8c — green
	Nord15 = RGB(180, 142, 173) // #b48ead — purple

	// Convenience aliases
	NordBg         = Nord0
	NordBgElevated = RGB(39, 44, 54) // #272c36
	NordBgScrim    = RGB(30, 33, 40) // #1e2128
)
