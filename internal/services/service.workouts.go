package services

import (
	"backend/internal/datasources/records"
	"backend/internal/datasources/repositories"
	"backend/internal/http/data_transfers"
	"backend/pkg/convert"
	"backend/third_party/io"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"net/http"
)

type WorkoutsRepository interface {
	FindAll() ([]records.Workouts, error)
	FindByID(id int) (records.Workouts, error)
	FindAllByOwnerID(ownerID int) ([]records.Workouts, error)
	FindAllByCurrentUserID(ownerID int) ([]records.Workouts, error)
	Save(workout records.Workouts) (int, error)
	Update(id int, workout map[string]interface{}) error
	Delete(id int) error
	Copy(id int, userID int) (int, error)
	FindAllWithFilters(params repositories.QueryParams) ([]records.Workouts, int, error)
	LikeWorkout(id int, userID int) error
}

type WorkoutsService struct {
	repository              WorkoutsRepository
	workoutExercisesService *WorkoutExercisesService
	exercisesService        *ExercisesService
	ionet                   *io.Client
}

func NewWorkoutsService(repository WorkoutsRepository, workoutExercisesService *WorkoutExercisesService, ionet *io.Client, exercisesService *ExercisesService) *WorkoutsService {
	return &WorkoutsService{
		repository:              repository,
		workoutExercisesService: workoutExercisesService,
		ionet:                   ionet,
		exercisesService:        exercisesService,
	}
}

func (s *WorkoutsService) FindAll() ([]data_transfers.WorkoutsResponse, int, error) {
	var workoutsResponse []data_transfers.WorkoutsResponse

	workouts, err := s.repository.FindAll()
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return nil, http.StatusNotFound, errors.New("workouts not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	err = copier.Copy(&workoutsResponse, &workouts)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for i, workout := range workoutsResponse {
		if workout.Price == float64(0) {
			workoutExercises, statusCode, err := s.workoutExercisesService.FindAllByWorkoutID(workout.ID)
			if err != nil {
				return nil, statusCode, err
			}
			workoutsResponse[i].Exercises = workoutExercises
		}
	}

	return workoutsResponse, http.StatusOK, nil
}

func (s *WorkoutsService) FindByID(id int) (data_transfers.WorkoutsResponse, int, error) {
	var workoutResponse data_transfers.WorkoutsResponse

	workout, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return workoutResponse, http.StatusNotFound, errors.New("workout not found")
		}
		return workoutResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&workoutResponse, &workout)
	if err != nil {
		return workoutResponse, http.StatusInternalServerError, err
	}

	if workout.Price == float64(0) {
		workoutExercises, statusCode, err := s.workoutExercisesService.FindAllByWorkoutID(workout.ID)
		if err != nil {
			return workoutResponse, statusCode, err
		}
		workoutResponse.Exercises = workoutExercises
	}

	return workoutResponse, http.StatusOK, nil
}

func (s *WorkoutsService) FindByIDByOwner(id int, userID int) (data_transfers.WorkoutsResponse, int, error) {
	var workoutResponse data_transfers.WorkoutsResponse

	workout, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return workoutResponse, http.StatusNotFound, errors.New("workout not found")
		}
		return workoutResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&workoutResponse, &workout)
	if err != nil {
		return workoutResponse, http.StatusInternalServerError, err
	}

	if workout.Price == float64(0) || workout.OwnerID == userID {
		workoutExercises, statusCode, err := s.workoutExercisesService.FindAllByWorkoutID(workout.ID)
		if err != nil {
			return workoutResponse, statusCode, err
		}
		workoutResponse.Exercises = workoutExercises
	}

	return workoutResponse, http.StatusOK, nil
}

func (s *WorkoutsService) FindAllByOwnerID(ownerID int) ([]data_transfers.WorkoutsResponse, int, error) {
	var workoutsResponse []data_transfers.WorkoutsResponse

	workouts, err := s.repository.FindAllByOwnerID(ownerID)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return nil, http.StatusNotFound, errors.New("workouts not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	var publicWorkouts []records.Workouts
	for _, workout := range workouts {
		if workout.IsPrivate == false {
			publicWorkouts = append(publicWorkouts, workout)
		}
	}

	err = copier.Copy(&workoutsResponse, &publicWorkouts)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for i, workout := range workoutsResponse {
		if workout.Price == float64(0) {
			workoutExercises, statusCode, err := s.workoutExercisesService.FindAllByWorkoutID(workout.ID)
			if err != nil {
				return nil, statusCode, err
			}
			workoutsResponse[i].Exercises = workoutExercises
		}
	}

	return workoutsResponse, http.StatusOK, nil
}

func (s *WorkoutsService) FindAllByCurrentUserID(ownerID int) ([]data_transfers.WorkoutsResponse, int, error) {
	var workoutsResponse []data_transfers.WorkoutsResponse

	workouts, err := s.repository.FindAllByOwnerID(ownerID)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return nil, http.StatusNotFound, errors.New("workouts not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	err = copier.Copy(&workoutsResponse, &workouts)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for i, workout := range workoutsResponse {
		workoutExercises, statusCode, err := s.workoutExercisesService.FindAllByWorkoutID(workout.ID)
		if err != nil {
			return nil, statusCode, err
		}
		workoutsResponse[i].Exercises = workoutExercises
	}

	return workoutsResponse, http.StatusOK, nil
}

func (s *WorkoutsService) Save(workout data_transfers.CreateWorkoutRequest) (int, int, error) {
	var workoutRecord records.Workouts
	err := copier.Copy(&workoutRecord, &workout)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	id, err := s.repository.Save(workoutRecord)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusCreated, nil
}

func (s *WorkoutsService) Update(id int, workout data_transfers.UpdateWorkoutRequest) (int, error) {
	workoutMap, err := convert.StructToMap(workout)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Update - convert.StructToMap: %w", err)
	}

	err = s.repository.Update(id, workoutMap)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *WorkoutsService) Delete(id int) (int, error) {
	err := s.repository.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *WorkoutsService) Copy(id int, userID int) (int, int, error) {
	id, err := s.repository.Copy(id, userID)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusCreated, nil
}

func (s *WorkoutsService) PurchaseWorkout(workoutID int, userID int) (int, int, error) {
	// TODO: purchase workout

	id, err := s.repository.Copy(workoutID, userID)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusCreated, nil
}

func (s *WorkoutsService) FindAllWithFilters(params repositories.QueryParams) ([]data_transfers.WorkoutsResponse, int, int, error) {
	var workoutsResponse []data_transfers.WorkoutsResponse

	workouts, total, err := s.repository.FindAllWithFilters(params)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return nil, 0, http.StatusNotFound, errors.New("workouts not found")
		}
		return nil, 0, http.StatusInternalServerError, err
	}

	var freeWorkouts []records.Workouts
	for _, workout := range workouts {
		if workout.Price == float64(0) {
			freeWorkouts = append(freeWorkouts, workout)
		}
	}

	err = copier.Copy(&workoutsResponse, &freeWorkouts)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, err
	}

	for i, workout := range workoutsResponse {
		workoutExercises, statusCode, err := s.workoutExercisesService.FindAllByWorkoutID(workout.ID)
		if err != nil {
			return nil, 0, statusCode, err
		}
		workoutsResponse[i].Exercises = workoutExercises

	}

	return workoutsResponse, total, http.StatusOK, nil
}

func (s *WorkoutsService) LikeWorkout(id int, userID int) (int, error) {
	err := s.repository.LikeWorkout(id, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return http.StatusNotFound, errors.New("workout not found")
		}
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *WorkoutsService) GenerateWorkout(generateRequest data_transfers.WorkoutGenerateRequest) (int, error) {

	response, err := s.ionet.GenerateWorkout(generateRequest)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - GenerateWorkout - ionet.GenerateWorkout: %w", err)
	}
	workout := data_transfers.CreateWorkoutRequest{
		Title:       response.WorkoutName,
		Description: response.Description,
		IsPrivate:   false,
		Price:       0,
		OwnerID:     generateRequest.OwnerID,
	}
	workoutID, statusCode, err := s.Save(workout)
	if err != nil {
		return statusCode, fmt.Errorf("service - GenerateWorkout - Save: %w", err)
	}

	fmt.Println(response.Exercises)
	for _, exercise := range response.Exercises {
		exercise2, statusCode, err := s.exercisesService.FindByName(exercise.Name)
		if err != nil || exercise2.ID == 0 {
			fmt.Println("Exercise not found, creating new one", exercise.Name)
			customExercise := data_transfers.CreateExercisesRequest{
				Name: exercise.Name,
			}
			exerciseID, statusCode, err := s.exercisesService.CreateCustomExercise(customExercise, generateRequest.OwnerID)
			if err != nil {
				return statusCode, fmt.Errorf("service - GenerateWorkout - CreateCustomExercise: %w", err)
			}
			workoutExercise := data_transfers.CreateWorkoutExercisesRequest{
				WorkoutID:  workoutID,
				ExerciseID: exerciseID,
				MainNote:   exercise.Description,
				OwnerID:    generateRequest.OwnerID,
			}

			_, statusCode, err = s.workoutExercisesService.Save(workoutExercise)
			if err != nil {
				return statusCode, fmt.Errorf("service - GenerateWorkout - Save: %w", err)

			}

			fmt.Println("Exercise created", customExercise.Name)
		} else {
			fmt.Println("Exercise found", exercise.Name)
			workoutExercise := data_transfers.CreateWorkoutExercisesRequest{
				WorkoutID:  workoutID,
				ExerciseID: exercise2.ID,
				MainNote:   exercise.Description,
				OwnerID:    generateRequest.OwnerID,
			}
			_, statusCode, err = s.workoutExercisesService.Save(workoutExercise)
			if err != nil {
				return statusCode, fmt.Errorf("service - GenerateWorkout - Save: %w", err)
			}
			fmt.Println("Exercise added to workout")
		}

	}

	return http.StatusCreated, nil
}
