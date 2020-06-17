#! /bin/bash
FAAS_SHA=3e4762e025c570fe66d549af243800c79cd3630e
WATCHDOG_SHA=c796e1b714d703c90bf6e3f392471746e0aeab2d

(
	rm -rf build/
	mkdir -p build/
)

(
	cd build/

	wget https://github.com/openfaas/faas/archive/${FAAS_SHA}.zip
	wget https://github.com/openfaas-incubator/of-watchdog/archive/${WATCHDOG_SHA}.zip

	unzip ${FAAS_SHA}.zip -d .
	mv faas-${FAAS_SHA}/ faas/

	unzip ${WATCHDOG_SHA}.zip -d .
	mv of-watchdog-${WATCHDOG_SHA}/ faas/of-watchdog/

	rm ${FAAS_SHA}.zip
	rm ${WATCHDOG_SHA}.zip

	(
		cd faas/

		patch -p1 < ../../gateway.patch
		patch -p1 < ../../of-watchdog.patch
		patch -p1 < ../../build.patch
	)

	cp -R ../accesscontrol/ faas/accesscontrol/
	cp -R ../sample-hello-retail/ faas/sample-hello-retail/
)
