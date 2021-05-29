#version 410

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vp;

out vec3 color;

void main() {
  gl_Position = projection * camera * model * vec4(vp, 1.0);
  color = 1 / vp;
}
