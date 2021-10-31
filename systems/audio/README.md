# Audio preface

How will this work, create a system that contains the audio context and Entity
handle mechanisms, when an entity is added the system will check if it is has
the specific components
i.e:

- "AudioListener" to be placed in cameras
- "AudioSource" to be placed in whatever emits a sound

## TODO

Create a positional audio processor that redirects volume to specific channels
i.e:

- output is 2 channels do left and right + some pass filter on back
- output is 4 channels redirect volume based on orientation and distance
