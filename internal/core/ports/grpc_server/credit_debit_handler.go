package grpc_server

import (
	"context"

	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CreditDebitNoteHandler struct {
	pb.UnimplementedCreditDebitNoteServiceServer
	svc ports.CreditDebitNoteService
}

func NewGRPCServer(svc ports.CreditDebitNoteService) *CreditDebitNoteHandler {
	return &CreditDebitNoteHandler{svc: svc}
}

func (s *CreditDebitNoteHandler) CreateCreditDebitNote(ctx context.Context, req *pb.CreateCreditDebitNoteRequest) (*pb.CreditDebitNote, error) {
    noteReq := req.Note
    if noteReq == nil {
        return nil, status.Errorf(codes.InvalidArgument, "note is required")
    }

    invoiceID, err := uuid.Parse(noteReq.InvoiceId)
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid invoice_id UUID format: %v", err)
    }

    note := domain.CreditDebitNote{
        InvoiceID: invoiceID,
        Type:      domain.NoteType(noteReq.Type.String()), 
        Amount:    noteReq.Amount.String(),               
        CreatedBy: noteReq.Audit.GetCreatedBy(),          
    }

    createdNote, err := s.svc.Create(ctx, note)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create note: %v", err)
    }

    return mapDomainToProtoCreditDebitNote(createdNote), nil
}


func (s *CreditDebitNoteHandler) GetCreditDebitNote(ctx context.Context, req *pb.GetCreditDebitNoteRequest) (*pb.CreditDebitNote, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id UUID format: %v", err)
	}

	note, err := s.svc.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "note not found: %v", err)
	}

	return mapDomainToProtoCreditDebitNote(note), nil
}

func (s *CreditDebitNoteHandler) ListCreditDebitNotes(ctx context.Context, req *pb.ListCreditDebitNotesRequest) (*pb.ListCreditDebitNotesResponse, error) {
    page := req.GetPage()
    var limit int32 = 50  // sensible default
    var offset int32 = 0

    if page != nil {
        limit = page.PageSize
    }

    notes, err := s.svc.List(ctx, limit, offset)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to list notes: %v", err)
    }

    protoNotes := make([]*pb.CreditDebitNote, len(notes))
    for i, n := range notes {
        protoNotes[i] = mapDomainToProtoCreditDebitNote(n)
    }

    // TODO: generate a next_page_token if you want proper pagination
    return &pb.ListCreditDebitNotesResponse{
        Notes: protoNotes,
        Page: &pb.PageResponse{
            NextPageToken: "",
            TotalSize:     int64(len(protoNotes)), // or actual count from DB
        },
    }, nil
}


func (s *CreditDebitNoteHandler) UpdateCreditDebitNote(ctx context.Context, req *pb.UpdateCreditDebitNoteRequest) (*pb.CreditDebitNote, error) {
    noteReq := req.GetNote()
    if noteReq == nil {
        return nil, status.Errorf(codes.InvalidArgument, "note is required")
    }

    id, err := uuid.Parse(noteReq.Id)
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid id UUID format: %v", err)
    }

    invoiceID, err := uuid.Parse(noteReq.InvoiceId)
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid invoice_id UUID format: %v", err)
    }

    note := domain.CreditDebitNote{
        ID:        id,
        InvoiceID: invoiceID,
        Type:      domain.NoteType(noteReq.Type.String()), // map enum → domain
        Amount:    noteReq.Amount.String(),                // map Money → string, or convert properly
        Reason:    noteReq.Reason,
        UpdatedBy: noteReq.Audit.GetUpdatedBy(),           // comes from AuditFields inside CreditDebitNote
    }

    // Optionally: use req.UpdateMask to only update specific fields in DB
    updatedNote, err := s.svc.Update(ctx, note)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to update note: %v", err)
    }

    return mapDomainToProtoCreditDebitNote(updatedNote), nil
}

func (s *CreditDebitNoteHandler) DeleteCreditDebitNote(ctx context.Context, req *pb.DeleteCreditDebitNoteRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id UUID format: %v", err)
	}
	if err := s.svc.Delete(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete note: %v", err)
	}
	return &emptypb.Empty{}, nil
}
