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

func (t *Tool) Clone() Tool {
	na := []string{}
	na = append(na, t.Args...)

	return Tool{
		Command: t.Command,
		Args:    na,
	}
}

func (r Recipe) Clone() *Recipe {
	nr := Recipe{}
	for _, t := range r {
		nr = append(nr, t.Clone())
	}
	return &nr
}

func (r Recipe) Len() int {
	return len(r)
}

func (g *GenerateOptions) Clone() *GenerateOptions {
	o := []string{}
	o = append(o, g.Create...)

	return &GenerateOptions{
		Create: o,
	}
}

func (b *BuildOptions) Clone() *BuildOptions {
	return &BuildOptions{
		BuildRecipe: b.BuildRecipe.Clone(),
		Cleanup:     b.Cleanup.Clone(),
	}
}

func (gc *CleanupOptions) Clone() *CleanupOptions {
	var ng *CleanupGlobOptions
	var nc *CleanupCommandOptions

	if gc == nil {
		return nil
	}

	if gc.Command != nil {
		nc = gc.Command.Clone()
	}

	if gc.Glob != nil {
		ng = gc.Glob.Clone()
	}

	return &CleanupOptions{
		Glob:    ng,
		Command: nc,
	}
}

func (gc *CleanupGlobOptions) Clone() *CleanupGlobOptions {
	np := []string{}
	np = append(np, gc.Patterns...)

	return &CleanupGlobOptions{
		Recursive: gc.Recursive,
		Patterns:  np,
	}
}

func (gc *CleanupCommandOptions) Clone() *CleanupCommandOptions {
	return &CleanupCommandOptions{
		Recipe: gc.Recipe.Clone(),
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
