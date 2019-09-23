package network

import "fmt"
import "golang.org/x/sys/unix"
import "net"
import "log"
import "encoding/binary"
import "flag"

// CANMsg almacena un mensaje de la red CAN
type CANMsg struct {
	id  uint32
	dlc uint8
	ext uint8
	rtr uint8
	//padding uint8
	//res0 uint8
	//res1 uint8
	data []byte
}

func byteArrayToCANMsg(array []byte, canmsg *CANMsg) {
	canmsg.id = binary.LittleEndian.Uint32(array[0:4])

	if canmsg.id&unix.CAN_RTR_FLAG != 0 {
		canmsg.rtr = 1
	}
	if canmsg.id&unix.CAN_EFF_FLAG != 0 {
		canmsg.ext = 1
		canmsg.id = canmsg.id & unix.CAN_EFF_MASK
	} else {
		canmsg.ext = 0
		canmsg.id = canmsg.id & unix.CAN_SFF_MASK
	}

	canmsg.dlc = array[4]
	canmsg.data = array[8 : 8+array[4]]
}

func main() {
	canIface := flag.String("canif", "can0", "The CAN interface")

	flag.Parse()
	iface, err := net.InterfaceByName(*canIface)
	if err != nil {
		log.Fatal(err)
	}

	fd, _ := unix.Socket(unix.AF_CAN, unix.SOCK_RAW, unix.CAN_RAW)
	addr := &unix.SockaddrCAN{Ifindex: iface.Index}
	unix.Bind(fd, addr)
	frame := make([]byte, 16)
	// frame[0:3]: ID (LSB primero)
	// frame[4]: largo
	// frame[5]: padding
	// frame[6:7]: reserved
	// frame[8:15]: data

	canmsg := new(CANMsg)

	fmt.Println("Hola")
	for {
		unix.Read(fd, frame)

		byteArrayToCANMsg(frame, canmsg)
		fmt.Printf("%+v\n", *canmsg)
	}
}
