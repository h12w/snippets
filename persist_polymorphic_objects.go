package main

import (
    "encoding/json"
    "fmt"
    "reflect"
)

// This is a solution to persist polymorphic objects using json package as suggested by Andrew Gerrand
// on golang-nuts "Is there a way to persist polymorphic objects using xml/json package?"
type JsonObject struct {
    Type string
    Data json.RawMessage
}

func NewJsonObject(v interface{}) (*JsonObject, error) {
    b, err := json.Marshal(v)
    if err != nil {
        return nil, err
    }
    return &JsonObject{
        reflect.Indirect(reflect.ValueOf(v)).Type().Name(), b}, nil
}

type TypeMapper map[string]reflect.Type

func (m TypeMapper) GetObject(o *JsonObject) (interface{}, error) {
    t, found := m[o.Type]
    if !found {
        return nil, fmt.Errorf("Failed to find type: %s", o.Type)
    }

    v := reflect.New(t).Interface()
    err := json.Unmarshal(o.Data, v)
    if err != nil {
        return nil, err
    }
    return v, nil
}

// ---------------------------------------------------------

// ---------------------------------------------------------

type DrawingContext interface {
    // ...
}

type Shape interface {
    Draw(ctx DrawingContext)
}

type ShapeBase struct {
    X, Y, Width, Height float64
}

type Rectangle struct {
    ShapeBase ShapeBase
}

func (s *Rectangle) Draw(ctx DrawingContext) {
    // ...
}

type Eclipse struct {
    ShapeBase ShapeBase
}

func (s *Eclipse) Draw(ctx DrawingContext) {
    // ...
}

var ShapeMapper = TypeMapper{
    "Rectangle": reflect.TypeOf(Rectangle{}),
    "Eclipse":   reflect.TypeOf(Eclipse{}),
}

func main() {
    shapes := []Shape{
        &Rectangle{ShapeBase{5, 6, 7, 8}},
        &Eclipse{ShapeBase{1, 2, 3, 4}},
    }

    // marshal
    jsonObjects := []*JsonObject{}
    for _, shape := range shapes {
        jsonObject, _ := NewJsonObject(shape)
        jsonObjects = append(jsonObjects, jsonObject)
    }
    jsn, _ := json.MarshalIndent(jsonObjects, "", "  ")

    // unmarshal
    jsonObjects = []*JsonObject{}
    _ = json.Unmarshal(jsn, &jsonObjects)
    v := []Shape{}
    for _, jsonObject := range jsonObjects {
        obj, _ := ShapeMapper.GetObject(jsonObject)
        v = append(v, obj.(Shape))
    }

    fmt.Println(shapes[0], v[0])
    fmt.Println(shapes[1], v[1])
    fmt.Println(string(jsn))
}
