#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform vec2 time;
uniform vec2 resolution;
uniform vec2 movement_vector;

uniform float bullets_off_X;
uniform float bullets_off_Y;
uniform float border_size;

uniform float enemy_off_X;
uniform float powerup_off_X;

out vec4 finalColor;

vec3 col = vec3(0.0 , 0.0 , 0.0);



#define PI 3.14159265358979

vec3 color = 0.3*vec3(0.4,0.0,0.9);
float d2y(float d){return 1./(0.2+d);}
float radius = 0.3;

float fct(vec2 p, float r){
	float a = 3.*mod(atan(p.y, p.x)+time.x, 2.*PI);
	
	
	float scan = 20.*1.;
	return (d2y(a)+scan)*(1.-step(radius,r));
}

	
float circle(vec2 p, float r){
	float d=distance(r, radius);
	return d2y(10.*d);
}

void draw_star(vec2 uv, vec2 trans, float size) {
    col += size / length(uv - trans);
}

void draw_powerup(vec2 uv, vec2 trans, float size) {

    if (col.z > 0.5) {
        col += size / length( uv - trans);
    } else {
        col.z += size / length( uv - trans);
    }
}

void draw_red_star(vec2 uv, vec2 trans, float size) {
    col += size / length(uv - trans);
}

void draw_hole(vec2 uv) {
    col += 0.005 / (length(uv) - 0.2);
}

void draw_plane() {
    vec2 pos = ( gl_FragCoord.xy / resolution.xy ) - vec2(0.5,0.5);	
    float horizon = 0.0; 
    float fov = 0.5; 
	float scaling = time.y;
	
	vec3 p = vec3(pos.x, fov, pos.y - horizon);      
	vec2 s = vec2(p.x/p.z, p.y/p.z) * scaling;
		
	float run=(time.x/125)*0.2;
	if (pos.y>0.0)	
	  run=-run;
				
	float color = 1.-sign((mod(s.x*3. - run * movement_vector.x, 0.1) - 0.005 * length(pos)) * (mod(s.y*1. - run * movement_vector.y, 0.1) - 0.005 * length(pos)))*0.9+pos.y;		
	color *= p.z*p.z*10.0;

	if(pos.y<-0.05){
	    col += vec3(color) * 0.05;
	}
}

 
float snow(vec2 uv,float scale)
{
	float w=smoothstep(1.,0.,-uv.y*(scale/10.));
	if(w<.1)
		return 0.;
	uv+= (time.x/600)/scale;
	uv.y+=(time.x/600)*2./scale;
	uv.x+=sin(uv.y+(time.x/600)*.5)/scale;
	uv*=scale;
	vec2 s=floor(uv),f=fract(uv),p;
	float k=3.,d;
	p=.5+.35*sin(11.*fract(sin((s+p+scale)*mat2(7,3,6,5))*5.))-f;
	d=length(p);
	k=min(d,k);
	k=smoothstep(0.,k,sin(f.x+f.y)*0.02);
    	return k*w;
}
 
float rect( vec2 p, vec2 b)
{
	vec2 v = abs(p) - b;
	float d = length(max(v,0.0));
	return pow(d, 1.9);
}

void draw_snow(vec2 uv) {
    float c = snow(uv,10.);
	vec3 snow_c=(vec3(c));
    if(uv.y>-0.05){
	    col += vec3(snow_c);
	}
}

void draw_sun(vec2 uv, float size) {
	uv/=cos(.0*length(uv));
    uv /= size;
	float y  = 0.0;
	
	float dc = length(uv);
	
	y+=fct(uv, dc);
	y+=circle(uv, dc);
	//y+=grid(position, y);
	y=pow(y,1.80);

    if(uv.y>-0.05){
	    col += vec3( sqrt(y)*color);
	}
}

void draw_border(float size) {
    vec2 unipos = (gl_FragCoord.xy / resolution);
	vec2 pos = unipos*4.0-1.0;
		//pos.x *= resolution.x / resolution.y;	
	float d1 = rect(pos - vec2(1.0,1), vec2(0.5,1)); 
	vec3 clr1 = vec3(0.2,0.6,1.0) *d1; 
		
	col += vec3( clr1*size);
}

void main()
{   
    
    vec2 resolution = vec2(resolution.x, resolution.y);
    vec2 uv = (gl_FragCoord.xy - 0.5 *resolution.xy) / resolution.y;


    float sun_size = 2.5/(2-time.y);
    draw_snow(uv);
    draw_sun(uv, sun_size);

    draw_plane();

    if (time.y < 0.9) {
        draw_powerup(uv, vec2(powerup_off_X/25, 0), 0.01*(1-time.y));
        draw_red_star(uv, vec2(enemy_off_X/25, 0), 0.02*(1-time.y));
        draw_star(uv, vec2(bullets_off_X/25, (bullets_off_Y)/50 - 0.33), 1 / ((bullets_off_Y + 2) * 50) );
    }

    draw_border(border_size);
    
    finalColor = vec4(col, 1.0);
}