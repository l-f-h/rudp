package rudp

type packetListOrderType int8

const (
	packetListOrderBySeqNb packetListOrderType = iota
	packetListOrderByAckNb
)

// packetList is a sorted link list of packet that ordered by packet.SeqNumber/AckNumber ASC
type packetList struct {
	head, rail *node
	length     int32
	orderType  packetListOrderType
}

type node struct {
	data *packet
	next *node
}

func newPacketList() *packetList {
	return &packetList{}
}

func (l *packetList) getPacketSortKey(p *packet) uint32 {
	if l.orderType == packetListOrderBySeqNb {
		return p.seqNumber
	} else { // if l.orderType == packetListOrderByAckNb
		return p.ackNumber
	}
}

func (l *packetList) putPacket(p *packet) {
	newNode := &node{data: p}
	if l.head == nil {
		l.head = newNode
		l.rail = newNode
		l.length++
		return
	} else {
		var last *node
		cur := l.head
		sortKey := l.getPacketSortKey(p)
		for ; cur != nil; cur = cur.next {
			curSortKey := l.getPacketSortKey(cur.data)
			if curSortKey == sortKey {
				return
			} else if curSortKey > sortKey {
				newNode.next = cur
				if last == nil { // only one node now
					l.head = newNode
					l.rail = cur
				} else {
					last.next = newNode
				}
				l.length++
				return
			} else {
				last = cur
			}
		}
		l.rail.next = newNode
		l.rail = newNode
		l.length++
		return
	}
}

func (l *packetList) empty() bool {
	return l.length == 0
}

func (l *packetList) getHeadPacket() *packet {
	return l.head.data
}

// removePacketBySeqNb return val means if the seqNb packet is be found and deleted
func (l *packetList) removePacketByNb(nb uint32) {
	if l.empty() {
		return
	}

	var last *node
	for cur := l.head; cur != nil; cur = cur.next {
		curNb := l.getPacketSortKey(cur.data)
		if curNb == nb {
			l.length--
			if cur == l.head {
				if l.head == l.rail {
					l.head = nil
					l.rail = nil
				} else {
					l.head = l.head.next
				}
			} else {
				last.next = cur.next
				if cur == l.rail {
					l.rail = last
				}
			}
			return
		}
		last = cur
	}
}
