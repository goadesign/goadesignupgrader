package apidsl

func MediaType(identifier string, apidsl func()) interface{} {
	return nil
}

func Attribute(name string, args ...interface{}) {
	return
}

func Resource(name string, dsl func()) interface{} {
	return nil
}

func Action(name string, dsl func()) interface{} {
	return nil
}

func Routing(routs ...interface{}) interface{} {
	return nil
}

func GET(path string, dsl ...func()) interface{} {
	return nil
}

func BasePath(val string) {
	return
}

func HashOf(k, v interface{}, dsls ...func()) interface{} {
	return nil
}

func Metadata(name string, value ...string) {
	return
}

func Response(name string, paramsAndDSL ...interface{}) {
	return
}

func Status(status int) {
	return
}

func CollectionOf(v interface{}, paramAndDSL ...interface{}) interface{} {
	return nil
}

func Type(name string, apidsl func()) interface{} {
	return nil
}

func API(name string, dsl func()) interface{} {
	return nil
}

func Param(name string, args ...interface{}) {
	return
}

func Params(dsl func()) {
	return
}

func CanonicalActionName(a string) {
	return
}
