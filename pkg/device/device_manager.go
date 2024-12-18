package device

import (
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var logger = log.Default()

type DeviceManager struct {
	devicePath    string
	deviceDrivers []DeviceDriver
}

func NewDeviceManager(devicePath string) (*DeviceManager, error) {
	deviceManager := DeviceManager{
		devicePath: devicePath,
	}

	var err error
	deviceManager.deviceDrivers, err = deviceManager.GetAllDevices()

	if err != nil {
		return nil, err
	}

	return &deviceManager, nil
}

func (m DeviceManager) GetAllDevices() ([]DeviceDriver, error) {
	devicesFile, err := os.ReadDir(m.devicePath)
	if err != nil {
		log.Fatal(err)
	}

	devices := make([]DeviceDriver, 0)

	for _, deviceFile := range devicesFile {

		yamlFile, err := os.ReadFile(path.Join(m.devicePath, deviceFile.Name()))
		if err != nil {
			println("Failed -", deviceFile.Name())
			continue
		}

		var device Device
		err = yaml.Unmarshal(yamlFile, &device)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
			continue
		}

		switch device.Driver {
		case HeaterType:
			devices = append(devices, NewHeater(device))
		case InverterType:
			devices = append(devices, NewInterver(device))
		default:
			devices = append(devices, DeviceDriver(device))
		}

	}

	return devices, nil

}

func (m DeviceManager) GetDeviceByName(name string) (DeviceDriver, error) {
	for _, dev := range m.deviceDrivers {
		if dev.GetDeviceName() == name {
			return dev, nil
		}
	}

	return nil, errNotFound
}

func (m DeviceManager) GetDeviceDriver(id DriverType) (DeviceDriver, error) {
	for _, dev := range m.deviceDrivers {
		if dev.GetDriverName() == id {
			return dev, nil
		}
	}

	return nil, errNotFound
}

func (m DeviceManager) PowerOn(id DriverType) error {
	device, err := m.GetDeviceDriver(id)

	if err != nil {
		return err
	}

	return device.PowerOn()
}

func (m DeviceManager) PowerOff(id DriverType) error {
	device, err := m.GetDeviceDriver(id)

	if err != nil {
		return err
	}

	return device.PowerOff()
}

func (m DeviceManager) DeviceStatus(id DriverType) (DeviceStatus, error) {
	device, err := m.GetDeviceDriver(id)

	if err != nil {
		return INACTIVE, err
	}

	return device.Status()
}

// Status() (DeviceStatus, error)
// ReadValue() (any, error)
// GetConfig() (DeviceConfig, error)
