package auth

//ProfileIndex type profile index.
//You should use ProfileIndex(string) as profile index.
type ProfileIndex string

//ProfileIndexName profile index name
const ProfileIndexName = ProfileIndex("Name")

//ProfileIndexFirstName profile first name
const ProfileIndexFirstName = ProfileIndex("FirstName")

//ProfileIndexMiddleName profile last name
const ProfileIndexMiddleName = ProfileIndex("LastName")

//ProfileIndexLastName profile last name
const ProfileIndexLastName = ProfileIndex("LastName")

//ProfileIndexEmail profile index email
const ProfileIndexEmail = ProfileIndex("Email")

//ProfileIndexNickname profile index nick name
const ProfileIndexNickname = ProfileIndex("Nickname")

//ProfileIndexAvatar profile index avatar
const ProfileIndexAvatar = ProfileIndex("Avatar")

//ProfileIndexProfileURL profile index url
const ProfileIndexProfileURL = ProfileIndex("ProfileURL")

//ProfileIndexAccessToken profile index accesstoken
const ProfileIndexAccessToken = ProfileIndex("AccessToken")

//ProfileIndexGender profile index Gender
const ProfileIndexGender = ProfileIndex("Gender")

//ProfileIndexCompany profile index company
const ProfileIndexCompany = ProfileIndex("Company")

//ProfileIndexID profile index id.
const ProfileIndexID = ProfileIndex("ID")

//ProfileIndexLocation profile index location
const ProfileIndexLocation = ProfileIndex("Location")

//ProfileIndexWebsite profile index website
const ProfileIndexWebsite = ProfileIndex("Website")

//ProfileIndexLocale profile index locale
const ProfileIndexLocale = ProfileIndex("Locale")

//ProfileGenderMale Male value for profile gender field.
const ProfileGenderMale = "M"

//ProfileGenderFemale Female value for profile gender field.
const ProfileGenderFemale = "F"

//ProfileGenderUnknow Unknow value for profile gender field.
const ProfileGenderUnknow = ""

//ProfileIndexCountry profile index country
const ProfileIndexCountry = ProfileIndex("Country")

//ProfileIndexProvince profile index province
const ProfileIndexProvince = ProfileIndex("Province")

//ProfileIndexCity profile index city
const ProfileIndexCity = ProfileIndex("City")

//Profile type stores user profile data
type Profile map[ProfileIndex][]string

//Value return first data of profile field.
//if profile is empty,empty string will be returned.
func (p *Profile) Value(index ProfileIndex) string {
	data, ok := (*p)[index]
	if ok == false || len(data) == 0 {
		return ""
	}
	return data[0]
}

//Values return all data of profile field.
//if profile is empty,nil will be returned.
func (p *Profile) Values(index ProfileIndex) []string {
	data, ok := (*p)[index]
	if ok == false {
		return nil
	}
	return data
}

//SetValue set string as profile data and clear previous data.
func (p *Profile) SetValue(index ProfileIndex, value string) {
	(*p)[index] = []string{value}
}

//SetValues set string slice as profile data and clear previous data.
func (p *Profile) SetValues(index ProfileIndex, values []string) {
	(*p)[index] = values
}

//AddValue add value to profile field.
func (p *Profile) AddValue(index ProfileIndex, value string) {
	data, ok := (*p)[index]
	if ok == false {
		data = []string{}
	}
	data = append(data, value)
	(*p)[index] = data
}
