struct resolution
  width: int
  height: int
end

struct video_mode
  resolution: resolution
  interlaced: bool
  frame_rate: float
  name: box<string?>
end

def new_resolution(width, height)
  { resolution: width, height }
end

let res1 = { resolution: width = 0, height = 0 }
let res2 = new_resolution(0, 0)

if res1 == res2 then
  printf("res1 is equal to res2\n")
end

let video_mode = {
  video_mode: resolution = res1,
  interlaced = false,
  frame_rate = 0.0,
  name = box(none)
}

video_mode.name <- "noninterlaced video"
printf("video mode name: %s\n", unbox(video_mode.name))
