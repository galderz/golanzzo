package main

import (
	"fmt"
)

const CacheServiceFixedMemoryXmxMb = 200
const CacheServiceJvmNativeMb = 220
const CacheServiceMaxRamMb = CacheServiceFixedMemoryXmxMb + CacheServiceJvmNativeMb
const CacheServiceAdditionalJavaOptions = "-Dsun.zip.disableMemoryMapping=true -XX:+UseSerialGC -XX:MinHeapFreeRatio=5 -XX:MaxHeapFreeRatio=10"

func main() {
	extraJavaOptions := "-Xlog:gc*:stdout:time"

	javaOptions := fmt.Sprintf(
		"-Xmx%dM -Xms%dM -XX:MaxRAM=%dM %s %s",
		CacheServiceFixedMemoryXmxMb,
		CacheServiceFixedMemoryXmxMb,
		CacheServiceMaxRamMb,
		CacheServiceAdditionalJavaOptions,
		extraJavaOptions,
	)

	fmt.Println("javaOptions: ", javaOptions)
}
