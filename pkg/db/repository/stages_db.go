package db

import (
	"appContract/pkg/db"
	"appContract/pkg/models"
	"log"
)

func DBgetStageAll() ([]models.Stages, error ){
    conn, err:=db.ConnectDB()
    if err!=nil{
        log.Fatal(err)
    }
    defer conn.Close()

    rows, err:=conn.Query(`SELECT 
    s.id_stage,
    s.name_stage,
    s.id_user,
    u.surname,
    u.username,
    u.patronymic,
    u.phone,
    u.email,
    s.description,
    s.id_status_stage,
    ss.name_status_stage,
    s.date_create_start,
    s.date_create_end,
    s.id_contract,
    c.name_contract,
    c.date_create_contract
FROM stages s
JOIN users u ON s.id_user = u.id_user
JOIN contracts c ON s.id_contract = c.id_contract
JOIN status_stages ss ON s.id_status_stage = ss.id_status_stage`)

    if err!=nil{
        log.Fatal(err)      
    }
    defer rows.Close()

    var stages []models.Stages
    for rows.Next(){
        var stage models.Stages
        err=rows.Scan(&stage.Id_stage,
                      &stage.Name_stage,
                      &stage.Id_user,
                      &stage.Surname,
                      &stage.Username,
                      &stage.Patronymic,
                      &stage.Phone,
                      &stage.Email,
                      &stage.Description,
                      &stage.Id_status_stage,
                      &stage.Name_status_stage,
                      &stage.Data_create_start,
                      &stage.Date_create_end,
                      &stage.Id_contract,
                      &stage.Name_contract,
                      &stage.Data_contract_create)
        if err!=nil{
            log.Fatal(err)
        }
        stages=append(stages, stage)
    }
    return stages, nil
}

func DBgetStageUserID(user_id int) ([]models.Stages, error ){
    conn, err:=db.ConnectDB()
    if err!=nil{
        log.Fatal(err)
    }
    defer conn.Close()

    rows, err:=conn.Query(`SELECT 
    s.id_stage,
    s.name_stage,
    s.id_user,
    s.description,
    s.id_status_stage,
    ss.name_status_stage,
    s.date_create_start,
    s.date_create_end,
    s.id_contract,
    c.name_contract,
    c.date_create_contract
FROM stages s
JOIN contracts c ON s.id_contract = c.id_contract
JOIN status_stages ss ON s.id_status_stage = ss.id_status_stage 
WHERE s.id_user=$1`,user_id)
    if err!=nil{
        log.Fatal(err)      
    }
    defer rows.Close()

    var stages []models.Stages
    for rows.Next(){
        var stage models.Stages
        err=rows.Scan(&stage.Id_stage,
            &stage.Name_stage,
            &stage.Id_user,
            &stage.Description,
            &stage.Id_status_stage,
            &stage.Name_status_stage,
            &stage.Data_create_start,
            &stage.Date_create_end,
            &stage.Id_contract,
            &stage.Name_contract,
            &stage.Data_contract_create)
        if err!=nil{
            log.Fatal(err)
        }
        stages=append(stages, stage)
    }
    return stages, nil
}

func DBgetStageID(stage_id int) (models.Stages, error ){
    conn, err:=db.ConnectDB()
    if err!=nil{
        log.Fatal(err)
    }
    defer conn.Close()

    rows, err:=conn.Query(`SELECT 
    s.id_stage,
    s.name_stage,
    s.id_user,
    u.surname,
    u.username,
    u.patronymic,
    u.phone,
    u.email,
    s.description,
    s.id_status_stage,
    ss.name_status_stage,
    s.date_create_start,
    s.date_create_end,
    s.id_contract,
    c.name_contract,
    c.date_create_contract,
    t.id_type_contract,
    t.id_type_contract
FROM stages s
JOIN users u ON s.id_user = u.id_user
JOIN contracts c ON s.id_contract = c.id_contract
JOIN types_contracts t ON c.id_type = t.id_type_contract
JOIN status_stages ss ON s.id_status_stage = ss.id_status_stage 
WHERE s.id_stage=$1`,stage_id)
    if err!=nil{
        log.Fatal(err)      
    }
    defer rows.Close()   

    var stage models.Stages             
    for rows.Next(){
        err=rows.Scan(&stage.Id_stage,
                      &stage.Name_stage,
                      &stage.Id_user,
                      &stage.Surname,
                      &stage.Username,
                      &stage.Patronymic,
                      &stage.Phone,
                      &stage.Email,
                      &stage.Description,
                      &stage.Id_status_stage,
                      &stage.Name_status_stage,
                      &stage.Data_create_start,
                      &stage.Date_create_end,
                      &stage.Id_contract,
                      &stage.Name_contract,
                      &stage.Data_contract_create,
                      &stage.Id_type_contract,
                      &stage.Name_type_contract)  
        if err!=nil{
            log.Fatal(err)
        }
    }
    return stage, nil
}

func DBgetFilesStageID(id_files int) (models.File, error ){
    conn, err:=db.ConnectDB()
    if err!=nil{
        log.Fatal(err)
    }
    defer conn.Close()

    rows, err:=conn.Query("SELECT * FROM files WHERE id_files=$1",id_files)
    if err!=nil{
        log.Fatal(err)      
    }
    defer rows.Close()                  

    var file models.File
    for rows.Next(){
        err=rows.Scan(&file.Id_file,
                      &file.Name_file,
                      &file.Data,
                      &file.Type_file,
                      &file.Id_stage)
        if err!=nil{
            log.Fatal(err)
        }
    }
    return file, nil
}

func DBgetStageIdStatus(id_stage int) (models.Stages, error ){
    conn, err:=db.ConnectDB()
    if err!=nil{
        log.Fatal(err)
    }
    defer conn.Close()

    rows, err:=conn.Query(`SELECT * 
                            FROM status_stages
                            WHERE id_status_stage=$1`,id_stage)
    if err!=nil{
        log.Fatal(err)      
    }
    defer rows.Close()                  

    var stage models.Stages
    for rows.Next(){
        err=rows.Scan(&stage.Id_status_stage,
                      &stage.Name_status_stage)
        if err!=nil{
            log.Fatal(err)
        }
    }
    return stage, nil
}

    func DBaddFile(id_stage int, file models.File) error{
        conn, err:=db.ConnectDB()
        if err!=nil{
            log.Fatal(err)
        }
        defer conn.Close()
    
        _, err=conn.Exec(`INSERT INTO files 
        name_file,
        data,
        type_file,
        id_stage
        VALUES ($1,$2,$3,$4)`,
        file.Name_file,
        file.Data,
        file.Type_file,
        id_stage)
        
        if err!=nil{    
            log.Fatal(err)  
        }
        return nil
    }
func DBaddStage(stage models.Stages) error{
    conn, err:=db.ConnectDB()
    if err!=nil{
        log.Fatal(err)
    }
    defer conn.Close()

    _, err=conn.Exec(`INSERT INTO stages
    name_stage,
    id_user,
    description,
    id_status_stage,
    date_create_start,
    date_create_end,
    id_contract
    VALUES ($1,$2,$3,$4,$5,$6)`,
    stage.Name_stage,
    stage.Id_user,
    stage.Description,
    stage.Id_status_stage,
    stage.Data_create_start,
    stage.Date_create_end,
    stage.Id_contract)
    
    if err!=nil{    
        log.Fatal(err)  
    }
    return nil

}

func DBCreateComment(comment models.HistoreStatus) error{
    conn, err:=db.ConnectDB()
    if err!=nil{
        log.Fatal(err)
    }
    defer conn.Close()

    _, err=conn.Exec(`INSERT INTO histore_status
    id_history_state,
    id_status_stage,    
    id_stage,
    data_create,
    comment
    VALUES ($1,$2,$3,$4,$5)`,
    comment.Id_history_state,
    comment.Id_status_stage,
    comment.Id_stage,
    comment.Data_create,
    comment.Comment)
    
    if err!=nil{    
        log.Fatal(err)  
    }
    return nil

}

func DBgetComment(id_stage int) (models.HistoreStatus, error ){
    conn, err:=db.ConnectDB()
    if err!=nil{
        log.Fatal(err)
    }
    defer conn.Close()

    rows, err:=conn.Query(`SELECT * 
                            FROM history_states 
                            WHERE id_stage=$1`,id_stage)
    if err!=nil{
        log.Fatal(err)      
    }
    defer rows.Close()                  

    var comment models.HistoreStatus  
    for rows.Next(){
        err=rows.Scan(&comment.Id_history_state,
                      &comment.Id_status_stage, 
                      &comment.Id_stage,
                      &comment.Data_create,
                      &comment.Comment)

        if err!=nil{
            log.Fatal(err)
        }
    }
    return comment, nil
}

func DBChengeStatusStage(id_stage int, id_status_stage int, comment string) error{
    conn, err:=db.ConnectDB()
    if err!=nil{
        log.Fatal(err)
    }
    defer conn.Close()

    _, err=conn.Exec(`UPDATE stages 
    SET id_status_stage=$1, data_create=NOW()
    WHERE id_stage=$2`,

    id_status_stage,
    id_stage)
    
    if err!=nil{    
        log.Fatal(err)  
    }

    _, err=conn.Exec(`INSERT INTO history_states (id_stage, id_status_stage, data_create, comment)
    VALUES ($1, $2, NOW(), $3)`,

    id_stage,
    id_status_stage,
    comment)
    
    if err!=nil{    
        log.Fatal(err)  
    }
    return nil
}

func DBdeleteFile(id_files int) error{
    conn, err:=db.ConnectDB()
    if err!=nil{
        log.Fatal(err)
    }
    defer conn.Close()  

    _, err=conn.Exec(`DELETE FROM files WHERE id_files=$1`, id_files)   

    if err!=nil{
        log.Fatal(err)
    }
    return nil
}

func DBdeleteStage(id_stage int) error{
    conn, err:=db.ConnectDB()
    if err!=nil{
        log.Fatal(err)
    }
    defer conn.Close()

    _, err=conn.Exec(`DELETE FROM stages WHERE id_stage=$1`, id_stage)

    if err!=nil{
        log.Fatal(err)
    }
    return nil
}

