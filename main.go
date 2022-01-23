package main

import (
	"strconv";
	"math/rand";
	"github.com/gen2brain/raylib-go/raylib"
)

type bullet struct {
    x float32
    y float32
	size float32
}

type Player struct {
	score int
	powerup int
	x float32
	y float32

	speed float32
}

type Enemy struct {
	x float32
	y float32
}

type Powerup struct {
	x float32
	y float32
}

func main() {

	bullet := bullet{x:0, y:0}
	player := Player{score:0, powerup:0 , x:0, y:0, speed:1}
	enemy := Enemy{x:0, y:0}
	powerup := Powerup{x:0, y:0}

	screenWidth := int32(1366)
	screenHeight := int32(768)

	resolution := make([]float32, 2)
	resolution[0] = float32(screenWidth)
	resolution[1] = float32(screenHeight)

	icon := rl.LoadImage("data/icon.png")

	rl.InitWindow(screenWidth, screenHeight, "Abyss waft")
	rl.ToggleFullscreen();
	rl.SetWindowIcon(*icon);
	rl.SetExitKey(rl.KeyDelete);

	rl.InitAudioDevice()

	theme := rl.LoadMusicStream("data/sounds/theme.mp3")
	rl.SetMusicVolume(theme, 0.1)
	rl.PlayMusicStream(theme)

	click_sound := rl.LoadSound("data/sounds/click.wav")
	destroy_sound := rl.LoadSound("data/sounds/destroy.wav")
	powerup_sound := rl.LoadSound("data/sounds/powerup.wav")
	shoot_sound := rl.LoadSound("data/sounds/shoot.wav")
	rl.SetSoundVolume(click_sound, 0.05)
	rl.SetSoundVolume(destroy_sound, 0.05)
	rl.SetSoundVolume(powerup_sound, 0.05)
	rl.SetSoundVolume(shoot_sound, 0.05)


	//Shaders
	shader := rl.LoadShader("data/shaders/base_vert.vs", "data/shaders/base_frag.fs")
	default_sh := rl.LoadShader("data/shaders/base_vert.vs", "data/shaders/bases_frag.fs")
	glow := rl.LoadShader("sdata/haders/base_vert.vs", "data/shaders/glow.fs")

	ship := rl.LoadTexture("data/ship.png")

	timeUnicodeLoc := rl.GetShaderLocation(shader, "time")
	ResolutionUnicodeLoc := rl.GetShaderLocation(shader, "resolution")
	rl.SetShaderValue(shader, ResolutionUnicodeLoc, resolution, rl.ShaderUniformVec2)


	X_bullet_offset_UnicodeLoc := rl.GetShaderLocation(shader, "bullets_off_X")
	Y_bullet_offset_UnicodeLoc := rl.GetShaderLocation(shader, "bullets_off_Y")

	opacity_UnicodeLoc := rl.GetShaderLocation(default_sh, "opacity")

	
	border_size_UnicodeLoc := rl.GetShaderLocation(shader, "border_size")
	X_enemy_offset_UnicodeLoc := rl.GetShaderLocation(shader, "enemy_off_X")
	X_powerup_offset_UnicodeLoc := rl.GetShaderLocation(shader, "powerup_off_X")

	default_sh_ResolutionUnicodeLoc := rl.GetShaderLocation(default_sh, "resolution")
	rl.SetShaderValue(default_sh, default_sh_ResolutionUnicodeLoc, resolution, rl.ShaderUniformVec2)

	time := make([]float32, 2)
	time[0] = float32(5)
	time[1] = float32(0)

	movement_vector := make([]float32, 2)
	movement_vector[0] = float32(0)
	movement_vector[1] = float32(3)
	MovementVectorUnicodeLoc := rl.GetShaderLocation(shader, "movement_vector")

	border_size := make([]float32, 1)
	border_size[0] = 0.2
	rl.SetShaderValue(shader, border_size_UnicodeLoc, border_size , rl.ShaderUniformFloat)

	rl.SetTargetFPS(60)

	playing := false
	time[1] = 0.99
	shoot_sound_cd := 0
	for !rl.WindowShouldClose() {
		rl.UpdateMusicStream(theme)

		border_size[0] = ((float32(player.powerup)/17)+0.2)
		rl.SetShaderValue(shader, border_size_UnicodeLoc, border_size , rl.ShaderUniformFloat)

		if (playing == true) {
			if (time[1] > 0.1) {
				time[1] -= 0.01
			}
		} else {

			if (time[1] < 0.99) {
				time[1] += 0.01
			}
		}

		
		if (playing == false) {
			if rl.IsKeyDown(rl.KeySpace) {
				playing = true
				rl.PlaySoundMulti(click_sound)
			}
		}

		if rl.IsKeyDown(rl.KeyEscape) {
			playing = false
			player.score = 0
			player.powerup = 0
			player.speed = 1
			powerup.x = (-player.x/5) + 1600
			rl.PlaySoundMulti(click_sound)
		}

		time[0]++
		player.y++;		
		bullet.y += 1 * (1);

		if (playing == true) { 
			if rl.IsKeyDown(rl.KeyD) {
				if (movement_vector[0] < 2) {
					movement_vector[0] += float32(0.02)
				}
				player.x -= player.speed
			}

			if rl.IsKeyDown(rl.KeyA) {

				if (movement_vector[0] > -2) {
					movement_vector[0] -= float32(0.02)
				}
				player.x += player.speed
			}

			if rl.IsKeyDown(rl.KeySpace) {
				shoot_sound_cd++;
				if (shoot_sound_cd > 10) {
					rl.PlaySoundMulti(shoot_sound)
					shoot_sound_cd = 0
				}
				
				bullet.x = (-player.x/5)
				bullet.y = 0
			}

			if (bullet.y > 16) {
				if (bullet.x - 1 < enemy.x && bullet.x + 1 > enemy.x ) {
					player.score++;
					bullet.x = (-player.x/5) + 800
					bullet.y = 0
					rl.PlaySoundMulti(destroy_sound)
					shoud_i_spawn_powerup := rand.Intn(6)

					if (shoud_i_spawn_powerup >= 5) {
						powerup.x = enemy.x
						powerup.y = 0
					} else {
						powerup.x = (-player.x/5) + 1600
					}
			
					rand_e_of := float32(rand.Intn(25)) - 12.5
					enemy.x += rand_e_of
				}

				if (bullet.y > 17) { 
					bullet.x = (-player.x/5) + 800
					bullet.y = 0
				}
			}

			if (bullet.y > 16) {
				if (bullet.x - 1 < powerup.x && bullet.x + 1 > powerup.x ) {
					player.powerup++;
					rl.PlaySoundMulti(powerup_sound)
					player.speed += 0.1;
					bullet.x = (-player.x/5) + 800
					bullet.y = 0
					powerup.x = (-player.x/5) + 1600
				}
			}

			if (bullet.y > 17) { 
				bullet.x = (-player.x/5) + 800
				bullet.y = 0
			}
		}

		value_powerup_x := make([]float32, 1)
		value_powerup_x[0] = powerup.x + (player.x/5)
		rl.SetShaderValue(shader, X_powerup_offset_UnicodeLoc, value_powerup_x , rl.ShaderUniformFloat)


		value_enemy_x := make([]float32, 1)
		value_enemy_x[0] = enemy.x + (player.x/5)
		rl.SetShaderValue(shader, X_enemy_offset_UnicodeLoc, value_enemy_x , rl.ShaderUniformFloat)

		value_y := make([]float32, 1)
		value_y[0] = bullet.y
		rl.SetShaderValue(shader, Y_bullet_offset_UnicodeLoc, value_y , rl.ShaderUniformFloat)

		value_x := make([]float32, 1)
		value_x[0] = bullet.x + (player.x/5)
		rl.SetShaderValue(shader, X_bullet_offset_UnicodeLoc, value_x , rl.ShaderUniformFloat)

		rl.SetShaderValue(shader, timeUnicodeLoc, time, rl.ShaderUniformVec2)
		rl.SetShaderValue(shader, MovementVectorUnicodeLoc, movement_vector, rl.ShaderUniformVec2)

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		rl.BeginShaderMode(shader)
		rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.White)
		rl.EndShaderMode()

		
		if (playing == true) {
			ship_opacity := make([]float32, 1)
			ship_opacity[0] = 1 - time[1]
			rl.SetShaderValue(default_sh, opacity_UnicodeLoc, ship_opacity , rl.ShaderUniformFloat)

			rl.BeginShaderMode(default_sh)
			rl.DrawTexture(ship, (screenWidth/2) - (ship.Width/2), (screenHeight/2 + 300) - (ship.Height/2), rl.White) 
			rl.EndShaderMode()
		
		
			rl.BeginShaderMode(glow)
			rl.DrawText("Score "+strconv.Itoa(player.score), 0 + 20, 0 + 20, 48, rl.White)
			rl.DrawText("Powerups  "+strconv.Itoa(player.powerup), 0 + 20, 0 + 68, 28, rl.White)
			rl.EndShaderMode()
		} else {
			rl.BeginShaderMode(glow)
			rl.DrawText("Abyss waft", 688 - 350, 384 - 150, 128, rl.White)
			rl.DrawText("press space to begin", 688 - 250, 384 + 40, 48, rl.White)


			rl.DrawText("Music by Krabas", 10, 768 - 48, 32, rl.White)
			rl.DrawText("Code by Maksasj", 1366 - 285, 768 - 48, 32, rl.White)
			rl.EndShaderMode()
		}
		rl.EndDrawing()
	}
	rl.StopSoundMulti()
	rl.UnloadSound(click_sound) 
	rl.UnloadSound(destroy_sound) 
	rl.UnloadSound(powerup_sound) 
	rl.UnloadSound(shoot_sound) 

	rl.UnloadTexture(ship)
	rl.UnloadShader(shader)
	rl.CloseWindow()
}