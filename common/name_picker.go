package common

import "math/rand"

var NAMES = []string{
	"Aaron",
	"Abigail",
	"Adam",
	"Aiden",
	"Aisha",
	"Albert",
	"Alex",
	"Alexander",
	"Alfie",
	"Alice",
	"Amber",
	"Amelia",
	"Amelie",
	"Amy",
	"Anna",
	"Annabelle",
	"Archie",
	"Arthur",
	"Austin",
	"Ava",
	"Beatrice",
	"Bella",
	"Benjamin",
	"Bethany",
	"Blake",
	"Bobby",
	"Brooke",
	"Caleb",
	"Callum",
	"Cameron",
	"Charles",
	"Charlie",
	"Charlotte",
	"Chloe",
	"Connor",
	"Daisy",
	"Daniel",
	"Darcey",
	"Darcy",
	"David",
	"Dexter",
	"Dylan",
	"Edward",
	"Eleanor",
	"Elijah",
	"Eliza",
	"Elizabeth",
	"Ella",
	"Ellie",
	"Elliot",
	"Elliott",
	"Elsie",
	"Emilia",
	"Emily",
	"Emma",
	"Erin",
	"Esme",
	"Ethan",
	"Eva",
	"Evan",
	"Evelyn",
	"Evie",
	"Faith",
	"Felix",
	"Finlay",
	"Finley",
	"Florence",
	"Francesca",
	"Frankie",
	"Freddie",
	"Frederick",
	"Freya",
	"Gabriel",
	"George",
	"Georgia",
	"Grace",
	"Gracie",
	"Hannah",
	"Harley",
	"Harriet",
	"Harrison",
	"Harry",
	"Harvey",
	"Heidi",
	"Henry",
	"Hollie",
	"Holly",
	"Hugo",
	"Ibrahim",
	"Imogen",
	"Isaac",
	"Isabel",
	"Isabella",
	"Isabelle",
	"Isla",
	"Isobel",
	"Ivy",
	"Jack",
	"Jacob",
	"Jake",
	"James",
	"Jamie",
	"Jasmine",
	"Jayden",
	"Jenson",
	"Jessica",
	"Joseph",
	"Joshua",
	"Jude",
	"Julia",
	"Kai",
	"Katie",
	"Kian",
	"Lacey",
	"Layla",
	"Leah",
	"Leo",
	"Leon",
	"Lewis",
	"Lexi",
	"Liam",
	"Lilly",
	"Lily",
	"Logan",
	"Lola",
	"Louie",
	"Louis",
	"Luca",
	"Lucas",
	"Lucy",
	"Luke",
	"Lydia",
	"Maddison",
	"Madison",
	"Maisie",
	"Maria",
	"Martha",
	"Maryam",
	"Mason",
	"Matilda",
	"Matthew",
	"Max",
	"Maya",
	"Megan",
	"Mia",
	"Michael",
	"Millie",
	"Mohammad",
	"Mohammed",
	"Mollie",
	"Molly",
	"Muhammad",
	"Nathan",
	"Niamh",
	"Noah",
	"Oliver",
	"Olivia",
	"Ollie",
	"Oscar",
	"Owen",
	"Paige",
	"Phoebe",
	"Poppy",
	"Reuben",
	"Riley",
	"Robert",
	"Ronnie",
	"Rory",
	"Rose",
	"Rosie",
	"Ruby",
	"Ryan",
	"Samuel",
	"Sara",
	"Sarah",
	"Scarlett",
	"Sebastian",
	"Seth",
	"Sienna",
	"Skye",
	"Sofia",
	"Sonny",
	"Sophia",
	"Sophie",
	"Stanley",
	"Summer",
	"Teddy",
	"Theo",
	"Theodore",
	"Thomas",
	"Tilly",
	"Toby",
	"Tommy",
	"Tyler",
	"Victoria",
	"Violet",
	"William",
	"Willow",
	"Zachary",
	"Zara",
	"Zoe",
}

func PickRandomName() string {
	return NAMES[rand.Intn(len(NAMES))]
}
