package ingress

// var scheme = runtime.NewScheme()
//
// func init() {
//	_ = netv1.AddToScheme(scheme)
// }
//
// func Load(path string) (*netv1.Ingress, error) {
//	data, err := os.ReadFile(path)
//	if err != nil {
//		return nil, err
//	}
//
//	decoder := serializer.NewCodecFactory(scheme).UniversalDeserializer()
//
//	obj, _, err := decoder.Decode(data, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	ing, ok := obj.(*netv1.Ingress)
//	if !ok {
//		return nil, fmt.Errorf("file does not contain a networking.k8s.io/v1 Ingress")
//	}
//
//	return ing, nil
// }
