package types

func (check *Checker) substituteTypes(context, typ Type, aliases *TypeAliases, seen map[Type]Type) Type {
	if typ == nil {
		return nil
	}
	if seen == nil {
		seen = make(map[Type]Type)
	}
	if seen[typ] != nil {
		return seen[typ]
	}
	seen[typ] = typ

	var sub Type

	switch t := typ.(type) {
	case *Array:
		sub = &Array{t.len, check.substituteTypes(context, t.elem, aliases, seen)}
	case *Slice:
		sub = &Slice{check.substituteTypes(context, t.elem, aliases, seen)}
	case *Struct:
		sub = &Struct{
			check.substituteTypesVars(context, t.fields, aliases, seen),
			// t.fields,
			t.tags,
			t.offsets,
			t.offsetsOnce,
			check.substituteTypesTypeNames(context, t.typeParams, aliases, seen),
			// t.typeParams,
		}

	case *Pointer:
		sub = &Pointer{check.substituteTypes(context, t.elem, aliases, seen)}
	case *Tuple:
		sub = &Tuple{check.substituteTypesVars(context, t.vars, aliases, seen)}
	case *Signature:
		sub = &Signature{
			t.scope,
			// check.substituteTypesVar(context, t.recv, aliases, seen),
			t.recv,
			check.substituteTypesTuple(context, t.params, aliases, seen),
			check.substituteTypesTuple(context, t.results, aliases, seen),
			t.variadic,
			check.substituteTypesTypeNames(context, t.typeParams, aliases, seen),
			// t.typeParams,
		}

	case *Interface:
		sub = &Interface{
			check.substituteTypesFuncs(context, t.methods, aliases, seen),
			// t.methods,
			check.substituteTypesNameds(context, t.embeddeds, aliases, seen),
			// t.embeddeds,
			check.substituteTypesFuncs(context, t.allMethods, aliases, seen),
			// t.allMethods,
		}

	case *Map:
		sub = &Map{
			check.substituteTypes(context, t.key, aliases, seen),
			check.substituteTypes(context, t.elem, aliases, seen),
		}
	case *Chan:
		sub = &Chan{t.dir, check.substituteTypes(context, t.elem, aliases, seen)}

	case *Named:
		sub = check.substituteTypesNamed(context, t, aliases, seen)
	default:
		sub = t
	}
	seen[typ] = sub
	return sub
}

func (check *Checker) substituteTypesNamed(context Type, old *Named, aliases *TypeAliases, seen map[Type]Type) Type {
	if old == nil {
		return nil
	}
	if aliases != nil && old.obj != nil && (*aliases)[old.obj] != nil && old.context == context {
		return (*aliases)[old.obj]
	} else {
		return old
	}
	// } else if old.underlying == old {
	// 	return old
	// } else {
	// 	return &Named{
	// 		old.obj,
	// 		check.substituteTypes(context, old.underlying, aliases, seen),
	// 		old.methods,
	// 		old.context,
	// 		old.variance,
	// 	}
	// }
}

func (check *Checker) substituteTypesObject(context Type, old object, aliases *TypeAliases, seen map[Type]Type) object {
	return object{old.parent, old.pos, old.pkg, old.name, check.substituteTypes(context, old.typ, aliases, seen), old.order_, old.scopePos_}
}

func (check *Checker) substituteTypesVar(context Type, old *Var, aliases *TypeAliases, seen map[Type]Type) *Var {
	if old == nil {
		return nil
	}
	return &Var{check.substituteTypesObject(context, old.object, aliases, seen), old.anonymous, old.visited, old.isField, old.used}
}

func (check *Checker) substituteTypesFunc(context Type, old *Func, aliases *TypeAliases, seen map[Type]Type) *Func {
	if old == nil {
		return nil
	}
	return &Func{check.substituteTypesObject(context, old.object, aliases, seen)}
}

func (check *Checker) substituteTypesTypeName(context Type, old *TypeName, aliases *TypeAliases, seen map[Type]Type) *TypeName {
	if old == nil {
		return nil
	}
	return &TypeName{check.substituteTypesObject(context, old.object, aliases, seen)}
}

func (check *Checker) substituteTypesNameds(context Type, old []*Named, aliases *TypeAliases, seen map[Type]Type) []*Named {
	if old == nil {
		return nil
	}
	nameds := make([]*Named, len(old))
	for i, v := range old {
		nameds[i] = check.substituteTypesNamed(context, v, aliases, seen).(*Named)
	}
	return nameds
}

func (check *Checker) substituteTypesVars(context Type, old []*Var, aliases *TypeAliases, seen map[Type]Type) []*Var {
	if old == nil {
		return nil
	}
	vars := make([]*Var, len(old))
	for i, v := range old {
		vars[i] = check.substituteTypesVar(context, v, aliases, seen)
	}
	return vars
}

func (check *Checker) substituteTypesFuncs(context Type, old []*Func, aliases *TypeAliases, seen map[Type]Type) []*Func {
	if old == nil {
		return nil
	}
	funcs := make([]*Func, len(old))
	for i, f := range old {
		funcs[i] = check.substituteTypesFunc(context, f, aliases, seen)
	}
	return funcs
}

func (check *Checker) substituteTypesTypeNames(context Type, old []*TypeName, aliases *TypeAliases, seen map[Type]Type) []*TypeName {
	if old == nil {
		return nil
	}
	names := make([]*TypeName, len(old))
	for i, t := range old {
		names[i] = check.substituteTypesTypeName(context, t, aliases, seen)
	}
	return names
}

func (check *Checker) substituteTypesTuple(context Type, old *Tuple, aliases *TypeAliases, seen map[Type]Type) *Tuple {
	if old == nil {
		return nil
	}
	return &Tuple{check.substituteTypesVars(context, old.vars, aliases, seen)}
}
