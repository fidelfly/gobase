package routex

type RouteConfig struct {
	restricted bool
	audit      bool
	props      map[string]interface{}
}

type RouteProps struct {
	props map[string]interface{}
}

func (rp *RouteProps) Set(key string, val interface{}) error {
	if rp.props == nil {
		return ErrPropIsNil
	}
	if _, ok := rp.props[key]; ok {
		return ErrPropConflict
	}

	rp.props[key] = val

	return nil
}

func (rp *RouteProps) Get(key string) (val interface{}, ok bool) {
	if rp.props == nil {
		return nil, false
	}
	val, ok = rp.props[key]
	return
}

func (rc *RouteConfig) SetProps(key string, prop interface{}) {
	if rc.props == nil {
		rc.props = make(map[string]interface{})
	}
	rc.props[key] = prop
}

func (rc *RouteConfig) GetProps(key string) interface{} {
	if rc.props == nil {
		return nil
	}
	return rc.props[key]
}

func NewConfig() RouteConfig {
	return RouteConfig{audit: true}
}

func (rc RouteConfig) GetCopy(includeProp ...bool) RouteConfig {
	copyRc := RouteConfig{
		restricted: rc.restricted,
		audit:      rc.audit,
	}
	if len(includeProp) > 0 && includeProp[0] {
		copyRc.props = rc.props
	}
	return copyRc
}
func (rc RouteConfig) IsRestricted() bool {
	return rc.restricted
}

func (rc RouteConfig) IsAuditEnable() bool {
	return rc.audit
}
