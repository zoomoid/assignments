package config

func (c *ConfigurationSpec) Clone() *ConfigurationSpec {
	nm := c.Members
	if nm != nil {
		nm = []GroupMember{}
		for _, m := range c.Members {
			nm = append(nm, m.Clone())
		}
	}

	ni := c.Includes
	if ni != nil {
		ni = []Include{}
		for _, i := range c.Includes {
			ni = append(ni, i.Clone())
		}
	}

	generateOptions := c.GenerateOptions
	if generateOptions != nil {
		generateOptions = c.GenerateOptions.Clone()
	}

	buildOptions := c.BuildOptions
	if buildOptions != nil {
		buildOptions = c.BuildOptions.Clone()
	}

	bundleOptions := c.BundleOptions
	if bundleOptions != nil {
		bundleOptions = c.BundleOptions.Clone()
	}

	return &ConfigurationSpec{
		Course:          c.Course,
		Group:           c.Group,
		Template:        c.Template,
		Members:         nm,
		Includes:        ni,
		GenerateOptions: generateOptions,
		BuildOptions:    buildOptions,
		BundleOptions:   bundleOptions,
	}
}

func (c *ConfigurationStatus) Clone() *ConfigurationStatus {
	return &ConfigurationStatus{
		Assignment: c.Assignment,
	}
}

func (c *Configuration) Clone() *Configuration {
	return &Configuration{
		Spec:   c.Spec.Clone(),
		Status: c.Status.Clone(),
	}
}

func (g *GroupMember) Clone() GroupMember {
	return GroupMember{
		Name: g.Name,
		ID:   g.ID,
	}
}

func (b *Include) Clone() Include {
	return Include{
		Path: b.Path,
	}
}

func (r *Recipe) Clone() Recipe {
	na := []string{}
	na = append(na, r.Args...)

	return Recipe{
		Command: r.Command,
		Args:    na,
	}
}

func (g *GenerateOptions) Clone() *GenerateOptions {
	o := []string{}
	o = append(o, g.Create...)

	return &GenerateOptions{
		Create: o,
	}
}

func (b *BuildOptions) Clone() *BuildOptions {
	nr := []Recipe{}
	for _, r := range b.Recipe {
		nr = append(nr, r.Clone())
	}
	return &BuildOptions{
		Recipe: nr,
	}
}

func (b *BundleOptions) Clone() *BundleOptions {
	nm := make(map[string]interface{})
	for k, v := range b.Data {
		nm[k] = v
	}

	i := []string{}
	i = append(i, b.Include...)

	return &BundleOptions{
		Template: b.Template,
		Data:     nm,
		Include:  i,
	}
}
