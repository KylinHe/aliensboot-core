/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/10/25
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package config

type DBConfig struct {
	Name       string 	//database name
	Address    string 	//database address
	MaxSession uint   	//database connection pool limit
	DialTimeout int   	//database dial timeout (second)
	SocketTimeout int
	SyncTimeout int 	//database session timeout (second)
	Mode *int
	QueryLimit int     //database query result max limit
}
