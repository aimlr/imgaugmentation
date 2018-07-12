# Description
Augments base images with different kinds of text data and outputs multiple image files.  These images can further be passed through differet image filters, transforms, etc. to create many variations,

We needed it to create a lot of variations of images that we could feed into our image recognition pipeline.

# Usage
```
./exe -n 100 -t examples/input/variations.json examples/input/india-passport-first-page.json
```

This will create 100 [default is 10] outputs based on the config in the json file and create variations of those using the specs in variations.json


