# Overview

A tiny JSON lib to read and alter a JSON value. The data structure is lazy, it's parse-on-read so that you can replace the parser with a faster one if performance is critical, use method `Transform` to do it.
