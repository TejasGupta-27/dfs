package p2p

//Peer is a wrapper for remote connection that represents a node 
type Peer interface{

}


//Transport -For handling communication between the nodes in the network 
//Can be of form (TCP,websockets)
type Transport interface{
	ListenAndAccept() error

}