package godump

var colors_disabled bool

func DisableColors() {
	colors_disabled = true
}

type colorizer struct {
	disabled bool
}

func (c colorizer) str(v string) string {
	if c.disabled || colors_disabled {
		return v
	}
	return "\033[38;2;138;201;38m" + v + "\033[0m"
}

func (c colorizer) quote(v string) string {
	if c.disabled || colors_disabled {
		return v
	}
	return "\033[38;2;112;214;255m" + v + "\033[0m"
}

func (c colorizer) bool(v string) string {
	if c.disabled || colors_disabled {
		return v
	}
	return "\033[38;2;249;87;56m" + v + "\033[0m"
}

func (c colorizer) num(v string) string {
	if c.disabled || colors_disabled {
		return v
	}
	return "\033[38;2;0;168;232m" + v + "\033[0m"
}

func (c colorizer) vtype(v string) string {
	if c.disabled || colors_disabled {
		return v
	}
	return "\033[38;2;0;150;199m" + v + "\033[0m"
}

func (c colorizer) ptr(v string) string {
	if c.disabled || colors_disabled {
		return v
	}
	return "\033[38;2;78;205;196m" + v + "\033[0m"
}

func (c colorizer) fn(v string) string {
	if c.disabled || colors_disabled {
		return v
	}
	return "\033[38;2;148;0;211m" + v + "\033[0m"
}

func (c colorizer) strct(v string) string {
	if c.disabled || colors_disabled {
		return v
	}
	return "\033[38;2;211;211;211m" + v + "\033[0m"
}

func (c colorizer) strct_prv(v string) string {
	if c.disabled || colors_disabled {
		return v
	}
	return "\033[38;2;255;123;0m" + v + "\033[0m"
}
