#version 330

varying vec2 fragTexCoord;
varying vec4 fragColor;

uniform sampler2D texture0;
uniform vec4 colDiffuse;

uniform float opacity;

void main()
{
    vec4 source = texture2D(texture0, fragTexCoord);
    gl_FragColor = source*vec4(opacity);
}