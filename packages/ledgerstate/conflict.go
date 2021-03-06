package ledgerstate

import (
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/iotaledger/hive.go/byteutils"
	"github.com/iotaledger/hive.go/cerrors"
	"github.com/iotaledger/hive.go/crypto"
	"github.com/iotaledger/hive.go/generics/objectstorage"
	"github.com/iotaledger/hive.go/marshalutil"
	"github.com/iotaledger/hive.go/stringify"
	"github.com/iotaledger/hive.go/types"
	"github.com/mr-tron/base58"
)

// region ConflictID ///////////////////////////////////////////////////////////////////////////////////////////////////

// ConflictIDLength contains the amount of bytes that a marshaled version of the ConflictID contains.
const ConflictIDLength = OutputIDLength

// ConflictID is the data type that represents the identifier of a Conflict.
type ConflictID [ConflictIDLength]byte

// NewConflictID creates a new ConflictID from an OutputID.
func NewConflictID(outputID OutputID) (conflictID ConflictID) {
	copy(conflictID[:], outputID[:])

	return
}

// ConflictIDFromMarshalUtil unmarshals a ConflictID using a MarshalUtil (for easier unmarshaling).
func ConflictIDFromMarshalUtil(marshalUtil *marshalutil.MarshalUtil) (conflictID ConflictID, err error) {
	conflictIDBytes, err := marshalUtil.ReadBytes(ConflictIDLength)
	if err != nil {
		err = errors.Errorf("failed to parse ConflictID (%v): %w", err, cerrors.ErrParseBytesFailed)
		return
	}
	copy(conflictID[:], conflictIDBytes)

	return
}

// ConflictIDFromRandomness returns a random ConflictID which can for example be used for unit tests.
func ConflictIDFromRandomness() (conflictID ConflictID) {
	crypto.Randomness.Read(conflictID[:])

	return
}

// OutputID returns the OutputID that the ConflictID represents.
func (c ConflictID) OutputID() (outputID OutputID) {
	outputID, _, err := OutputIDFromBytes(c.Bytes())
	if err != nil {
		panic(err)
	}

	return
}

// Bytes returns a marshaled version of the ConflictID.
func (c ConflictID) Bytes() []byte {
	return c[:]
}

// Base58 returns a base58 encoded version of the ConflictID.
func (c ConflictID) Base58() string {
	return base58.Encode(c.Bytes())
}

// String returns a human readable version of the ConflictID.
func (c ConflictID) String() string {
	return "ConflictID(" + c.Base58() + ")"
}

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////

// region ConflictIDs //////////////////////////////////////////////////////////////////////////////////////////////////

// ConflictIDs represents a collection of ConflictIDs.
type ConflictIDs map[ConflictID]types.Empty

// NewConflictIDs creates a new collection of ConflictIDs from the given list of ConflictIDs.
func NewConflictIDs(optionalConflictIDs ...ConflictID) (conflictIDs ConflictIDs) {
	conflictIDs = make(ConflictIDs)
	for _, conflictID := range optionalConflictIDs {
		conflictIDs[conflictID] = types.Void
	}

	return
}

// ConflictIDsFromMarshalUtil unmarshals a collection of ConflictIDs using a MarshalUtil (for easier unmarshaling).
func ConflictIDsFromMarshalUtil(marshalUtil *marshalutil.MarshalUtil) (conflictIDs ConflictIDs, err error) {
	conflictIDsCount, err := marshalUtil.ReadUint64()
	if err != nil {
		err = errors.Errorf("failed to parse count of ConflictIDs (%v): %w", err, cerrors.ErrParseBytesFailed)
		return
	}

	conflictIDs = make(ConflictIDs)
	for i := uint64(0); i < conflictIDsCount; i++ {
		conflictID, conflictIDErr := ConflictIDFromMarshalUtil(marshalUtil)
		if conflictIDErr != nil {
			err = errors.Errorf("failed to parse ConflictID: %w", conflictIDErr)
			return
		}

		conflictIDs[conflictID] = types.Void
	}

	return
}

// Add adds a ConflictID to the collection and returns the collection to enable chaining.
func (c ConflictIDs) Add(conflictID ConflictID) ConflictIDs {
	c[conflictID] = types.Void

	return c
}

// Slice returns a slice of ConflictIDs.
func (c ConflictIDs) Slice() (list []ConflictID) {
	list = make([]ConflictID, 0, len(c))
	for conflictID := range c {
		list = append(list, conflictID)
	}

	return
}

// Bytes returns a marshaled version of the ConflictIDs.
func (c ConflictIDs) Bytes() []byte {
	marshalUtil := marshalutil.New(marshalutil.Int64Size + len(c)*ConflictIDLength)
	marshalUtil.WriteUint64(uint64(len(c)))
	for conflictID := range c {
		marshalUtil.WriteBytes(conflictID.Bytes())
	}

	return marshalUtil.Bytes()
}

// String returns a human readable version of the ConflictIDs.
func (c ConflictIDs) String() string {
	if len(c) == 0 {
		return "ConflictIDs{}"
	}

	result := "ConflictIDs{\n"
	for conflictID := range c {
		result += strings.Repeat(" ", stringify.INDENTATION_SIZE) + conflictID.String() + ",\n"
	}
	result += "}"

	return result
}

// Clone creates a copy of the ConflictIDs.
func (c ConflictIDs) Clone() (clonedConflictIDs ConflictIDs) {
	clonedConflictIDs = make(ConflictIDs)
	for conflictID := range c {
		clonedConflictIDs[conflictID] = types.Void
	}

	return
}

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////

// region Conflict /////////////////////////////////////////////////////////////////////////////////////////////////////

// Conflict represents a set of Branches that are conflicting with each other.
type Conflict struct {
	id               ConflictID
	memberCount      int
	memberCountMutex sync.RWMutex

	objectstorage.StorableObjectFlags
}

// NewConflict is the constructor for new Conflicts.
func NewConflict(conflictID ConflictID) *Conflict {
	return &Conflict{
		id: conflictID,
	}
}

// FromObjectStorage creates a Conflict from sequences of key and bytes.
func (c *Conflict) FromObjectStorage(key, bytes []byte) (conflict objectstorage.StorableObject, err error) {
	conflict, err = c.FromBytes(byteutils.ConcatBytes(key, bytes))
	if err != nil {
		err = errors.Errorf("failed to parse Conflict from bytes: %w", err)
	}
	return
}

// FromBytes unmarshals a Conflict from a sequence of bytes.
func (c *Conflict) FromBytes(bytes []byte) (conflict *Conflict, err error) {
	marshalUtil := marshalutil.New(bytes)
	if conflict, err = c.FromMarshalUtil(marshalUtil); err != nil {
		err = errors.Errorf("failed to parse Conflict from MarshalUtil: %w", err)
		return
	}

	return
}

// FromMarshalUtil unmarshals a Conflict using a MarshalUtil (for easier unmarshaling).
func (c *Conflict) FromMarshalUtil(marshalUtil *marshalutil.MarshalUtil) (conflict *Conflict, err error) {
	if conflict = c; conflict == nil {
		conflict = &Conflict{}
	}
	if conflict.id, err = ConflictIDFromMarshalUtil(marshalUtil); err != nil {
		err = errors.Errorf("failed to parse ConflictID from MarshalUtil: %w", err)
		return
	}
	memberCount, err := marshalUtil.ReadUint64()
	if err != nil {
		err = errors.Errorf("failed to parse member count (%v): %w", err, cerrors.ErrParseBytesFailed)
		return
	}
	conflict.memberCount = int(memberCount)

	return
}

// ID returns the identifier of this Conflict.
func (c *Conflict) ID() ConflictID {
	return c.id
}

// MemberCount returns the amount of Branches that are part of this Conflict.
func (c *Conflict) MemberCount() int {
	c.memberCountMutex.RLock()
	defer c.memberCountMutex.RUnlock()

	return c.memberCount
}

// IncreaseMemberCount increase the MemberCount of this Conflict.
func (c *Conflict) IncreaseMemberCount(optionalDelta ...int) (newMemberCount int) {
	delta := 1
	if len(optionalDelta) >= 1 {
		delta = optionalDelta[0]
	}

	c.memberCountMutex.Lock()
	defer c.memberCountMutex.Unlock()

	c.memberCount += delta
	c.SetModified()
	newMemberCount = c.memberCount

	return c.memberCount
}

// DecreaseMemberCount decreases the MemberCount of this Conflict.
func (c *Conflict) DecreaseMemberCount(optionalDelta ...int) (newMemberCount int) {
	delta := 1
	if len(optionalDelta) >= 1 {
		delta = optionalDelta[0]
	}

	c.memberCountMutex.Lock()
	defer c.memberCountMutex.Unlock()

	c.memberCount -= delta
	c.SetModified()
	newMemberCount = c.memberCount

	return
}

// Bytes returns a marshaled version of the Conflict.
func (c *Conflict) Bytes() []byte {
	return byteutils.ConcatBytes(c.ObjectStorageKey(), c.ObjectStorageValue())
}

// String returns a human readable version of the Conflict.
func (c *Conflict) String() string {
	return stringify.Struct("Conflict",
		stringify.StructField("id", c.ID()),
		stringify.StructField("memberCount", c.MemberCount()),
	)
}

// ObjectStorageKey returns the key that is used to store the object in the database. It is required to match the
// StorableObject interface.
func (c *Conflict) ObjectStorageKey() []byte {
	return c.id.Bytes()
}

// ObjectStorageValue marshals the Conflict into a sequence of bytes. The ID is not serialized here as it is only used as
// a key in the ObjectStorage.
func (c *Conflict) ObjectStorageValue() []byte {
	return marshalutil.New(marshalutil.Uint64Size).
		WriteUint64(uint64(c.MemberCount())).
		Bytes()
}

// code contract (make sure the type implements all required methods)
var _ objectstorage.StorableObject = &Conflict{}

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////

// region ConflictMember ///////////////////////////////////////////////////////////////////////////////////////////////

// ConflictMemberKeyPartition defines the partition of the storage key of the ConflictMember model.
var ConflictMemberKeyPartition = objectstorage.PartitionKey(ConflictIDLength, BranchIDLength)

// ConflictMember represents the relationship between a Conflict and its Branches. Since an Output can have a
// potentially unbounded amount of conflicting Consumers, we store the membership of the Branches in the corresponding
// Conflicts as a separate k/v pair instead of a marshaled list of members inside the Branch.
type ConflictMember struct {
	conflictID ConflictID
	branchID   BranchID

	objectstorage.StorableObjectFlags
}

// NewConflictMember is the constructor of the ConflictMember reference.
func NewConflictMember(conflictID ConflictID, branchID BranchID) *ConflictMember {
	return &ConflictMember{
		conflictID: conflictID,
		branchID:   branchID,
	}
}

// FromObjectStorage creates an ConflictMember from sequences of key and bytes.
func (c *ConflictMember) FromObjectStorage(key, bytes []byte) (conflictMember objectstorage.StorableObject, err error) {
	conflictMember, err = c.FromBytes(byteutils.ConcatBytes(key, bytes))
	if err != nil {
		err = errors.Errorf("failed to parse ConflictMember from bytes: %w", err)
	}
	return
}

// FromBytes unmarshals a ConflictMember from a sequence of bytes.
func (c *ConflictMember) FromBytes(bytes []byte) (conflictMember *ConflictMember, err error) {
	marshalUtil := marshalutil.New(bytes)
	if conflictMember, err = c.FromMarshalUtil(marshalUtil); err != nil {
		err = errors.Errorf("failed to parse ConflictMember from MarshalUtil: %w", err)
		return
	}

	return
}

// FromMarshalUtil unmarshals an ConflictMember using a MarshalUtil (for easier unmarshaling).
func (c *ConflictMember) FromMarshalUtil(marshalUtil *marshalutil.MarshalUtil) (conflictMember *ConflictMember, err error) {
	if conflictMember = c; conflictMember == nil {
		conflictMember = &ConflictMember{}
	}
	if conflictMember.conflictID, err = ConflictIDFromMarshalUtil(marshalUtil); err != nil {
		err = errors.Errorf("failed to parse ConflictID from MarshalUtil: %w", err)
		return
	}
	if conflictMember.branchID, err = BranchIDFromMarshalUtil(marshalUtil); err != nil {
		err = errors.Errorf("failed to parse BranchID: %w", err)
		return
	}

	return
}

// ConflictID returns the identifier of the Conflict that this ConflictMember belongs to.
func (c *ConflictMember) ConflictID() ConflictID {
	return c.conflictID
}

// BranchID returns the identifier of the Branch that this ConflictMember references.
func (c *ConflictMember) BranchID() BranchID {
	return c.branchID
}

// Bytes returns a marshaled version of this ConflictMember.
func (c *ConflictMember) Bytes() []byte {
	return c.ObjectStorageKey()
}

// String returns a human readable version of this ConflictMember.
func (c *ConflictMember) String() string {
	return stringify.Struct("ConflictMember",
		stringify.StructField("conflictID", c.conflictID),
		stringify.StructField("branchID", c.branchID),
	)
}

// ObjectStorageKey returns the key that is used to store the object in the database. It is required to match the
// StorableObject interface.
func (c *ConflictMember) ObjectStorageKey() []byte {
	return byteutils.ConcatBytes(c.conflictID.Bytes(), c.branchID.Bytes())
}

// ObjectStorageValue marshals the Output into a sequence of bytes. The ID is not serialized here as it is only used as
// a key in the ObjectStorage.
func (c *ConflictMember) ObjectStorageValue() []byte {
	return nil
}

// code contract (make sure the type implements all required methods)
var _ objectstorage.StorableObject = &ConflictMember{}

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////
