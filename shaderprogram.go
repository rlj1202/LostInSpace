package lostinspace

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-compatibility/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ShaderProgram struct {
	program        uint32
	vertexShader   uint32
	fragmentShader uint32
}

func NewShaderProgram(vertexShaderRaw, fragmentShaderRaw string) *ShaderProgram {
	shaderProgram := &ShaderProgram{}

	program := gl.CreateProgram()
	shaderProgram.program = program
	vertexShader, err := compileShader(vertexShaderRaw, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderRaw, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	shaderProgram.vertexShader = vertexShader
	shaderProgram.fragmentShader = fragmentShader
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	return shaderProgram
}

func compileShader(rawSource string, shaderType uint32) (uint32, error) {
	rawSource += "\x00"
	shader := gl.CreateShader(shaderType)

	source, free := gl.Strs(rawSource)
	gl.ShaderSource(shader, 1, source, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("Failed to compile shader %v: %v\n", source, log)
	}

	return shader, nil
}

// Use this program
func (shaderProgram *ShaderProgram) Bind() {
	gl.UseProgram(shaderProgram.program)
}

// Get uniform location
func (shaderProgram *ShaderProgram) UniformLoc(name string) int32 {
	return gl.GetUniformLocation(shaderProgram.program, gl.Str(name+"\x00"))
}

// Set uniform value
func (shaderProgram *ShaderProgram) UniformInt(name string, value int32) {
	loc := shaderProgram.UniformLoc(name)
	gl.ProgramUniform1i(shaderProgram.program, loc, value)
}

// Set uniform value
func (shaderProgram *ShaderProgram) UniformMat4(name string, value mgl32.Mat4) {
	loc := shaderProgram.UniformLoc(name)
	gl.ProgramUniformMatrix4fv(shaderProgram.program, loc, 1, false, &(value[0]))
}
