-- Table events store all events
CREATE TABLE Events (
  id int NOT NULL AUTO_INCREMENT, -- unique id
  offset int NOT NULL, -- event offset
  type varchar(255) NOT NULL, -- event type ("EvX....")
  vc int NOT NULL, -- assigned vector clock
  ts bigint NOT NULL, -- timestamp
  g int NOT NULL, -- goroutine id
  p int NOT NULL, -- process id
  linkoff int, -- offset of the linked event
  predG int, -- predecessor goroutine
  predClk int, -- predecessor vector clock
  rid varchar(255), -- resource id (G, Ch, Wg, Cv, Mu)
  reid int, -- for channels only
  rval bigint, --  (for channels only) transmitted value
  rclock int, -- assigned vector clock to this particular rid
  stack_id int, -- stack id
  PRIMARY KEY (id)
);

-- Table stackFrames store the stack frame entries
CREATE TABLE StackFrames (
  id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  stack_id int NOT NULL,
  pc int NOT NULL,
  func varchar(255) NOT NULL,
  file varchar(255) NOT NULL,
  line int NOT NULL
);

-- Tabke args store argument values of the events that have arguments
CREATE TABLE Args (
  id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  eventID int NOT NULL,
  arg varchar(255) NOT NULL,
  value bigint NOT NULL
);

-- Table Goroutines maintain a model of goroutines tree
CREATE TABLE Goroutines (
  id int NOT NULL AUTO_INCREMENT,
  gid int NOT NULL,
  createStack_id int NOT NULL,
  parent_id int NOT NULL,
  ended int DEFAULT -1,
  PRIMARY KEY (id)
);
