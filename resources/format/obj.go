package format

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type OBJModel struct {
	Name         string
	MaterialFile string
	Groups       map[string]OBJGroup
}

type OBJGroup struct {
	UseMaterial []OBJUseMaterial
	Vertex      [][]float64
	UV          [][]float64
	Normal      [][]float64
	Polygon     [][]OBJPolygon
	Smoothing   int
}

type OBJUseMaterial struct {
	StartIndex uint64
	Material   string
}

type OBJPolygon struct {
	VertexIndex uint64
	UVIndex     uint64
	NormalIndex uint64
}

func LoadOBJ(reader io.Reader) (OBJModel, error) {
	result := OBJModel{
		Groups: map[string]OBJGroup{},
	}
	scanner := bufio.NewScanner(reader)
	group := ""
	for scanner.Scan() {
		line := scanner.Text()
		tokens := tokenize(line)
		if len(tokens) == 0 {
			continue
		}
		switch tokens[0] {
		case "mtllib":
			tokens = tokens[1:]
			if err := validLength(tokens, 1); err != nil {
				return OBJModel{}, err
			}
			result.MaterialFile = tokens[0]
		case "o":
			tokens = tokens[1:]
			if err := validLength(tokens, 1); err != nil {
				return OBJModel{}, err
			}
			result.Name = tokens[0]
		case "g":
			tokens = tokens[1:]
			if err := validLength(tokens, 1); err != nil {
				return OBJModel{}, err
			}
			group = tokens[0]
			result.Groups[group] = OBJGroup{}
		case "usemtl":
			tokens = tokens[1:]
			if err := validLength(tokens, 1); err != nil {
				return OBJModel{}, err
			}
			tmp := result.Groups[group]
			tmp.UseMaterial = append(tmp.UseMaterial, OBJUseMaterial{
				Material:   tokens[0],
				StartIndex: uint64(len(tmp.Polygon)),
			})
			result.Groups[group] = tmp
		case "v":
			tokens = tokens[1:]
			if err := validLength(tokens, 3); err != nil {
				return OBJModel{}, err
			}
			vertex, err := tokensToFloat64Values(tokens)
			if err != nil {
				return OBJModel{}, err
			}
			tmp := result.Groups[group]
			tmp.Vertex = append(tmp.Vertex, vertex)
			result.Groups[group] = tmp
		case "vt":
			tokens = tokens[1:]
			if err := validLength(tokens, 2); err != nil {
				return OBJModel{}, err
			}
			uv, err := tokensToFloat64Values(tokens)
			if err != nil {
				return OBJModel{}, err
			}
			tmp := result.Groups[group]
			tmp.UV = append(tmp.UV, uv)
			result.Groups[group] = tmp
		case "vn":
			tokens = tokens[1:]
			if err := validLength(tokens, 3); err != nil {
				return OBJModel{}, err
			}
			normal, err := tokensToFloat64Values(tokens)
			if err != nil {
				return OBJModel{}, err
			}
			tmp := result.Groups[group]
			tmp.Normal = append(tmp.Normal, normal)
			result.Groups[group] = tmp
		case "f":
			tokens = tokens[1:]
			if err := validLength(tokens, 3); err != nil {
				return OBJModel{}, err
			}
			polygon, err := tokensToOBJPolygons(tokens)
			if err != nil {
				return OBJModel{}, err
			}
			tmp := result.Groups[group]
			tmp.Polygon = append(tmp.Polygon, polygon)
			result.Groups[group] = tmp
		case "s":
			tokens = tokens[1:]
			if err := validLength(tokens, 1); err != nil {
				return OBJModel{}, err
			}
			smoothing, err := strconv.Atoi(tokens[0])
			if err != nil {
				return OBJModel{}, err
			}
			tmp := result.Groups[group]
			tmp.Smoothing = smoothing
			result.Groups[group] = tmp
		}
	}
	return result, nil
}

func tokensToOBJPolygons(tokens []string) ([]OBJPolygon, error) {
	result := make([]OBJPolygon, len(tokens))
	for i, s := range tokens {
		var (
			err     error
			polygon OBJPolygon
		)
		f := strings.Split(s, "/")
		if len(f) >= 1 && f[0] != "" {
			polygon.VertexIndex, err = strconv.ParseUint(f[0], 10, 64)
			if err != nil {
				return []OBJPolygon{}, err
			}
		}
		if len(f) >= 2 && f[1] != "" {
			polygon.UVIndex, err = strconv.ParseUint(f[1], 10, 64)
			if err != nil {
				return []OBJPolygon{}, err
			}
		}
		if len(f) >= 3 && f[2] != "" {
			polygon.NormalIndex, err = strconv.ParseUint(f[2], 10, 64)
			if err != nil {
				return []OBJPolygon{}, err
			}
		}
		result[i] = polygon
	}
	return result, nil
}

func tokensToFloat64Values(tokens []string) ([]float64, error) {
	result := make([]float64, len(tokens))
	for i, s := range tokens {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return []float64{}, err
		}
		result[i] = v
	}
	return result, nil
}

func validLength(src []string, length int) error {
	if len(src) != length {
		return errors.New("invalid format")
	}
	return nil
}

func tokenize(source string) []string {
	result := []string{}
	scanner := bufio.NewScanner(strings.NewReader(source))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result
}
