package server
type Tier string
const(TierFree Tier="free";TierPro Tier="pro")
type Limits struct{Tier Tier;Description string}
func LimitsFor(tier string)Limits{if tier=="pro"{return Limits{Tier:TierPro,Description:"Unlimited teams and members"}}
return Limits{Tier:TierFree,Description:"2 teams, 10 members"}}
func(l Limits)IsPro()bool{return l.Tier==TierPro}
