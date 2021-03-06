package dagsvisualizer

import "github.com/iotaledger/goshimmer/packages/jsonmodels"

const (
	// MsgTypeTangleVertex is the type of the Tangle DAG vertex.
	MsgTypeTangleVertex byte = iota
	// MsgTypeTangleBooked is the type of the Tangle DAG confirmed message.
	MsgTypeTangleBooked
	// MsgTypeTangleConfirmed is the type of the Tangle DAG confirmed message.
	MsgTypeTangleConfirmed
	// MsgTypeFutureMarkerUpdated is the type of the future marker updated message.
	MsgTypeFutureMarkerUpdated
	// MsgTypeUTXOVertex is the type of the UTXO DAG vertex.
	MsgTypeUTXOVertex
	// MsgTypeUTXOBooked is the type of the booked transaction.
	MsgTypeUTXOBooked
	// MsgTypeUTXOConfirmed is the type of the UTXO DAG vertex confirmed message.
	MsgTypeUTXOConfirmed
	// MsgTypeBranchVertex is the type of the branch DAG vertex.
	MsgTypeBranchVertex
	// MsgTypeBranchParentsUpdate is the type of the branch DAG vertex parents updated message.
	MsgTypeBranchParentsUpdate
	// MsgTypeBranchConfirmed is the type of the branch DAG vertex confirmed message.
	MsgTypeBranchConfirmed
	// MsgTypeBranchWeightChanged is the type of the branch DAG vertex weight changed message.
	MsgTypeBranchWeightChanged
)

type wsMessage struct {
	Type byte        `json:"type"`
	Data interface{} `json:"data"`
}

type tangleVertex struct {
	ID                      string   `json:"ID"`
	StrongParentIDs         []string `json:"strongParentIDs"`
	WeakParentIDs           []string `json:"weakParentIDs"`
	ShallowLikeParentIDs    []string `json:"shallowLikeParentIDs"`
	ShallowDislikeParentIDs []string `json:"shallowDislikeParentIDs"`
	BranchIDs               []string `json:"branchIDs"`
	IsMarker                bool     `json:"isMarker"`
	IsTx                    bool     `json:"isTx"`
	TxID                    string   `json:"txID,omitempty"`
	IsConfirmed             bool     `json:"isConfirmed"`
	ConfirmedTime           int64    `json:"confirmedTime"`
	GoF                     string   `json:"gof,omitempty"`
}

type tangleBooked struct {
	ID        string   `json:"ID"`
	IsMarker  bool     `json:"isMarker"`
	BranchIDs []string `json:"branchIDs"`
}

type tangleConfirmed struct {
	ID            string `json:"ID"`
	GoF           string `json:"gof"`
	ConfirmedTime int64  `json:"confirmedTime"`
}

type tangleFutureMarkerUpdated struct {
	ID             string `json:"ID"`
	FutureMarkerID string `json:"futureMarkerID"`
}

type utxoVertex struct {
	MsgID         string              `json:"msgID"`
	ID            string              `json:"ID"`
	Inputs        []*jsonmodels.Input `json:"inputs"`
	Outputs       []string            `json:"outputs"`
	IsConfirmed   bool                `json:"isConfirmed"`
	GoF           string              `json:"gof"`
	BranchIDs     []string            `json:"branchIDs"`
	ConfirmedTime int64               `json:"confirmedTime"`
}

type utxoBooked struct {
	ID        string   `json:"ID"`
	BranchIDs []string `json:"branchIDs"`
}

type utxoConfirmed struct {
	ID            string `json:"ID"`
	GoF           string `json:"gof"`
	ConfirmedTime int64  `json:"confirmedTime"`
}

type branchVertex struct {
	ID          string                                 `json:"ID"`
	Parents     []string                               `json:"parents"`
	IsConfirmed bool                                   `json:"isConfirmed"`
	Conflicts   *jsonmodels.GetBranchConflictsResponse `json:"conflicts"`
	GoF         string                                 `json:"gof"`
	AW          float64                                `json:"aw"`
}

type branchParentUpdate struct {
	ID      string   `json:"ID"`
	Parents []string `json:"parents"`
}

type branchConfirmed struct {
	ID string `json:"ID"`
}

type branchWeightChanged struct {
	ID     string  `json:"ID"`
	Weight float64 `json:"weight"`
	GoF    string  `json:"gof"`
}

type searchResult struct {
	Messages []*tangleVertex `json:"messages"`
	Txs      []*utxoVertex   `json:"txs"`
	Branches []*branchVertex `json:"branches"`
	Error    string          `json:"error,omitempty"`
}
