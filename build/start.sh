#!/bin/bash

set -ex -o pipefail

LANG=en_US.UTF-8
LANGUAGE=en_US:en
LC_ALL=en_US.UTF-8
ANDROID_HOME=/opt/sdk
PLATFORM=!platform!
EMULATOR_ROOT="${ANDROID_HOME}/emulator"
PATH="$PATH:${ANDROID_HOME}/tools:${ANDROID_HOME}/platform-tools:${ANDROID_HOME}/tools/bin:${ANDROID_HOME}/emulator"

# Export library paths
export ANDROID_SDK_ROOT=${ANDROID_HOME}
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:${EMULATOR_ROOT}/lib64/qt/lib:${EMULATOR_ROOT}/lib64/libstdc++:${EMULATOR_ROOT}/lib64:${EMULATOR_ROOT}/lib64/gles_swiftshader
export LIBGL_DEBUG=verbose

# Emulator options
CONSOLE_PORT=5554
ADB_PORT=5555

# Setup AVD
EMULATOR_CONFIG="$HOME/.android/avd/x86_64.avd/config.ini"
if [[ ! -d "$HOME/.android/avd/" ]] ; then
    echo | ${ANDROID_HOME}/tools/bin/avdmanager \
      create avd \
      --name "x86_64" \
      --package "system-images;${PLATFORM};google_apis;x86_64" \
      --tag google_apis

    # Fix the system image path (not sure why it breaks without this)
    sed -i 's:image.sysdir.1=sdk/:image.sysdir.1=:g' "${EMULATOR_CONFIG}"
    # Increase the resolution
    echo 'skin.name=480x800' >> "${EMULATOR_CONFIG}"
fi

# Set heap size to max if not defined in the environment
if [[ -z "${EMULATOR_HEAP_SIZE}" ]] ; then
    EMULATOR_HEAP_SIZE=576
fi
sed -i -E "s:(vm\.heapSize)=(.*):\1=${EMULATOR_HEAP_SIZE}:" "${EMULATOR_CONFIG}"

# Set default emulator options if not defined in the environment
if [[ -z "${EMULATOR_OPTS}" ]] ; then
    EMULATOR_OPTS="-screen multi-touch -no-boot-anim -noaudio -netfast -verbose -skip-adb-auth -no-snapshot -no-snapstorage"
fi
# Add any extra opts
EMULATOR_OPTS+=" ${EXTRA_EMULATOR_OPTS}"

# Set RAM configuration if provided in the environment
if [[ -z "${EMULATOR_RAM_SIZE}" ]] ; then
    EMULATOR_RAM_SIZE=4096
fi
sed -i -E "s:(hw\.ramSize)=(.*):\1=${EMULATOR_RAM_SIZE}:" "${EMULATOR_CONFIG}"

# Set CPU configuration
if [[ -z "${EMULATOR_NUM_CORE}" ]] ; then
    EMULATOR_NUM_CORE=2
fi
sed -i -E "s:(hw\.cpu\.ncore)=(.*):\1=${EMULATOR_NUM_CORE}:" "${EMULATOR_CONFIG}"

# Configure Play Store
if [[ "${EMULATOR_ENABLE_PLAYSTORE,,}" == "true" ]] ; then
    grep -q '^PlayStore' "${EMULATOR_CONFIG}" && \
    sed -i -E "s:(PlayStore\.enabled)=(.*):\1=true:" "${EMULATOR_CONFIG}" || echo 'PlayStore.enabled=true' >> "${EMULATOR_CONFIG}"
fi

# Remove any stale locks
find /root/.android/avd -name *lock -exec rm {} \;

# Start the emulator
${EMULATOR_ROOT}/qemu/linux-x86_64/qemu-system-x86_64-headless \
  -avd x86_64 \
  -gpu auto \
  -no-window \
  -ports ${CONSOLE_PORT},${ADB_PORT} \
  ${EMULATOR_OPTS} &

EMULATOR_PID=$!

adb wait-for-device
boot_completed=`adb -e shell getprop sys.boot_completed 2>&1`
timeout=0
until [ "X${boot_completed:0:1}" = "X1" ]; do
    sleep 5
    boot_completed=`adb shell getprop sys.boot_completed 2>&1 | head -n 1`
    echo "Read boot_completed property: <$boot_completed>"
    let "timeout += 1"
    if [ $timeout -gt 300 ]; then
         echo "Failed to start emulator"
         exit 1
    fi
done
adb wait-for-device

if [[ "$(ls /opt/sdk/apps/)" != "" ]]; then
  for i in /opt/sdk/apps/*.apk ; do
    pkg_search=$(basename $i | cut -d '.' -f1)
    adb shell "pm list packages" | grep "${pkg_search}" || adb install -r "$i"
  done
fi

wait ${EMULATOR_PID}
