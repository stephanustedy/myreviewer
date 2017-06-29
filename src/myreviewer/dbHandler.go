package myreviewer

import(
	"log"
)

func GetTeamList() ([]*Team, error) {
	teams := []*Team{}
	rows, err := db.Query(`
			SELECT
				team_id
				, name
				, coalesce(webhook, '')
			FROM
				team
			WHERE
				status = 1
		`)
	if err != nil {
		log.Println(err)
		return teams, err
	}
	defer rows.Close()
	for rows.Next() {
		team := &Team{}
		err := rows.Scan(&team.TeamId, &team.Name, &team.Webhook)
		if err != nil {
			log.Println(err)
			return teams, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func getMemberList() {

}

func deleteTeam() {

}

func deleteMember() {

}

func addMember() {

}

func addTeam(team *Team) error {
	stmt, err := db.Prepare("INSERT INTO team(name, status, webhook, channel) values(?,?,?,?)")
    if err != nil {
    	return err
    }

    _, err = stmt.Exec(team.Name, 1, team.Webhook, team.Channel)
    if err != nil {
    	return err
    }

    return nil
}

func updateTeam() {

}

func updateMember() {

}